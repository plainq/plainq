package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plainq/plainq/internal/houston"
	v1 "github.com/plainq/plainq/internal/server/schema/v1"
	"github.com/plainq/plainq/internal/server/storage"
	"github.com/plainq/servekit/errkit"
	"github.com/plainq/servekit/idkit"
	"github.com/plainq/servekit/respond"
)

func (s *PlainQ) createAccountHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name,omitempty"`
	}

	var req request

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("%w: decode request json: %s", errkit.ErrInvalidArgument, err.Error()))
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			s.logger.Error("create account: close request body",
				slog.String("error", err.Error()),
			)
		}
	}()

	if req.Name != "" {
		if err := validateUserName(req.Name); err != nil {
			respond.ErrorHTTP(w, r, fmt.Errorf("validate user name: %w", err))
			return
		}
	}

	if err := validatePassword(req.Password); err != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("validate user password: %w", err))
		return
	}

	hashedPassword, hashErr := s.hasher.HashPassword(req.Password)
	if hashErr != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("hash user password: %w", hashErr))
		return
	}

	verified := true

	if s.cfg.AuthRegistrationEnable {
		verified = false
	}

	userAccount := storage.Account{
		ID:        idkit.ULID(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		Verified:  verified,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.account.CreateAccount(r.Context(), userAccount); err != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("create user record: %w", hashErr))
		return
	}

	respond.Status(w, r, http.StatusCreated)
}

func (s *PlainQ) createQueueHandler(w http.ResponseWriter, r *http.Request) {
	var input v1.CreateQueueRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respond.ErrorHTTP(w, r, err)
		return
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			s.logger.Error("create queue: close request body",
				slog.String("error", err.Error()),
			)
		}
	}()

	output, createErr := s.queue.CreateQueue(r.Context(), &input)
	if createErr != nil {
		respond.ErrorHTTP(w, r, createErr)
		return
	}

	respond.JSON(w, r, output, respond.WithStatus(http.StatusCreated))
}

func (s *PlainQ) listQueuesHandler(w http.ResponseWriter, r *http.Request) {
	input := v1.ListQueuesRequest{
		QueuePrefix: r.URL.Query().Get("prefix"),
		Cursor:      r.URL.Query().Get("cursor"),
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		limit, parseErr := strconv.Atoi(l)
		if parseErr != nil {
			respond.ErrorHTTP(w, r, fmt.Errorf("%w: invalid limit", errkit.ErrInvalidArgument))
			return
		}

		if limit < 1 {
			respond.ErrorHTTP(w, r, fmt.Errorf("%w: invalid limit", errkit.ErrInvalidArgument))
			return
		}

		input.Limit = int32(limit)
	}

	output, listErr := s.queue.ListQueues(r.Context(), &input)
	if listErr != nil {
		respond.ErrorHTTP(w, r, listErr)
		return
	}

	respond.JSON(w, r, output)
}

func (s *PlainQ) describeQueueHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := validateQueueID(id); err != nil {
		respond.ErrorHTTP(w, r, err)
		return
	}

	input := v1.DescribeQueueRequest{QueueId: id}

	output, describeErr := s.queue.DescribeQueue(r.Context(), &input)
	if describeErr != nil {
		respond.ErrorHTTP(w, r, describeErr)
		return
	}

	respond.JSON(w, r, output, respond.WithStatus(http.StatusOK))
}

func (s *PlainQ) deleteQueueHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := validateQueueID(id); err != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("validation error: %w", err))
		return
	}

	force, parseErr := strconv.ParseBool(r.URL.Query().Get("force"))
	if parseErr != nil {
		force = false
	}

	input := v1.DeleteQueueRequest{
		QueueId: id,
		Force:   force,
	}

	output, deleteErr := s.queue.DeleteQueue(r.Context(), &input)
	if deleteErr != nil {
		respond.ErrorHTTP(w, r, deleteErr)
		return
	}

	respond.JSON(w, r, output, respond.WithStatus(http.StatusOK))
}

func (s *PlainQ) purgeQueueHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := validateQueueID(id); err != nil {
		respond.ErrorHTTP(w, r, fmt.Errorf("validation error: %w", err))
		return
	}

	output, purgeErr := s.queue.PurgeQueue(r.Context(), &v1.PurgeQueueRequest{
		QueueId: id,
	})
	if purgeErr != nil {
		respond.ErrorHTTP(w, r, purgeErr)
		return
	}

	respond.JSON(w, r, output, respond.WithStatus(http.StatusOK))
}

func (*PlainQ) houstonStaticHandler(w http.ResponseWriter, r *http.Request) {
	routeCtx := chi.RouteContext(r.Context())
	pathPrefix := strings.TrimSuffix(routeCtx.RoutePattern(), "/*")

	http.StripPrefix(pathPrefix, http.FileServerFS(houston.Bundle())).ServeHTTP(w, r)
}

func dropPolicyToString(policy v1.EvictionPolicy) string {
	switch policy {
	case v1.EvictionPolicy_EVICTION_POLICY_DROP:
		return "Drop Message"

	case v1.EvictionPolicy_EVICTION_POLICY_DEAD_LETTER:
		return "Dead Letter Queue"

	default:
		return ""
	}
}
