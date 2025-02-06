package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/heartwilltell/hc"
	"github.com/plainq/plainq/internal/houston"
	"github.com/plainq/plainq/internal/server/config"
	"github.com/plainq/plainq/internal/server/middleware"
	"github.com/plainq/plainq/internal/server/service/account"
	"github.com/plainq/plainq/internal/server/service/queue"
	"github.com/plainq/plainq/internal/server/service/telemetry"
	"github.com/plainq/servekit"
	"github.com/plainq/servekit/authkit/hashkit"
	"github.com/plainq/servekit/grpckit"
	"github.com/plainq/servekit/httpkit"
	_ "google.golang.org/grpc/encoding/proto"
)

// PlainQ represents plainq logic.
type PlainQ struct {
	cfg      *config.Config
	logger   *slog.Logger
	queue    *queue.Service
	account  *account.Service
	hasher   hashkit.Hasher
	observer telemetry.Observer
}

// NewServer returns a pointer to a new instance of the PlainQ.
func NewServer(
	cfg *config.Config,
	logger *slog.Logger,
	checker hc.HealthChecker,
	queue *queue.Service,
	account *account.Service,
) (*servekit.Server, error) {
	// Create a server which holds and serve all listeners.
	server := servekit.NewServer(logger)

	pq := PlainQ{
		cfg:      cfg,
		logger:   logger,
		queue:    queue,
		account:  account,
		observer: telemetry.NewObserver(),
	}

	// Create the HTTP listener.
	httpListener, httpListenerErr := listenerHTTP(cfg, logger, checker)
	if httpListenerErr != nil {
		return nil, httpListenerErr
	}

	// Initialize and mount the HTTP API routes.
	httpListener.MountGroup("/api", func(api chi.Router) {
		api.Use(middleware.Logging(logger))
		api.Use(cors.AllowAll().Handler)

		api.Route("/v1", func(v1 chi.Router) {
			if cfg.AuthEnable {
				v1.Route("/account", func(account chi.Router) {
					account.Mount("/", pq.account)
				})
			}

			// Queue related routes.
			v1.Route("/queue", func(queue chi.Router) {
				queue.Mount("/", pq.queue)
			})
		})
	})

	// Initialize and mount the Houston UI related routes.
	// There are routes responsible for static assets.
	httpListener.MountGroup("/", func(ui chi.Router) {
		// Static assets.
		ui.Get("/*", pq.houstonStaticHandler)
	})

	// Register the HTTP listener with a server.
	server.RegisterListener("HTTP", httpListener)

	grpcListener, grpcListenerErr := grpckit.NewListenerGRPC(cfg.GRPCAddr)
	if grpcListenerErr != nil {
		return nil, fmt.Errorf("create gRPC listener: %w", grpcListenerErr)
	}

	// Mount the queue gRPC routes to the gRPC listener.
	grpcListener.Mount(pq.queue)

	// Register the gRPC listener with a server.
	server.RegisterListener("GRPC", grpcListener)

	return server, nil
}

func (s *PlainQ) houstonStaticHandler(w http.ResponseWriter, r *http.Request) {
	routePattern := chi.RouteContext(r.Context()).RoutePattern()
	pathPrefix := strings.TrimSuffix(routePattern, "/*")

	s.logger.Debug("houston static handler",
		slog.String("path", r.URL.Path),
		slog.String("route_pattern", routePattern),
		slog.String("path_prefix", pathPrefix),
	)

	http.StripPrefix(pathPrefix, http.FileServerFS(houston.Bundle())).
		ServeHTTP(w, r)
}

func listenerHTTP(cfg *config.Config, logger *slog.Logger, checker hc.HealthChecker) (*httpkit.ListenerHTTP, error) {
	httpListenerOpts := httpkit.NewListenerOption[httpkit.ListenerConfig](
		httpkit.WithLogger(logger),
		httpkit.WithHTTPServerTimeouts(
			httpkit.HTTPServerReadHeaderTimeout(cfg.HTTPReadHeaderTimeout),
			httpkit.HTTPServerReadTimeout(cfg.HTTPReadTimeout),
			httpkit.HTTPServerWriteTimeout(cfg.HTTPWriteTimeout),
			httpkit.HTTPServerIdleTimeout(cfg.HTTPIdleTimeout),
		),
	)

	if cfg.HealthEnable {
		httpListenerOpts = append(httpListenerOpts, httpkit.WithHealthCheck(
			httpkit.HealthCheckRoute(cfg.HealthRoute),
			httpkit.HealthChecker(checker),
		))
	}

	if cfg.MetricsEnable {
		httpListenerOpts = append(httpListenerOpts, httpkit.WithMetrics(
			httpkit.MetricsRoute(cfg.MetricsRoute),
			httpkit.MetricsAccessLog(cfg.MetricsRouteLogs),
			httpkit.MetricsMetricsForEndpoint(cfg.MetricsRouteMetrics),
		))
	}

	httpListener, err := httpkit.NewListenerHTTP(cfg.HTTPAddr, httpListenerOpts...)
	if err != nil {
		return nil, fmt.Errorf("create HTTP listener: %w", err)
	}

	return httpListener, nil
}