package main

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/plainq/plainq/internal/server/config"
	"github.com/plainq/plainq/internal/server/middleware"
	"github.com/plainq/plainq/internal/server/service/account"
	"github.com/plainq/plainq/internal/server/service/onboarding"
	"github.com/plainq/plainq/internal/server/service/rbac"
	"github.com/plainq/servekit/authkit/hashkit"
	"github.com/plainq/servekit/authkit/jwtkit"
)

// This is an example of how to integrate the onboarding service
// into the main PlainQ server. This should be adapted to your
// actual server initialization code.

func main() {
	// Example configuration - adapt to your actual config loading
	cfg := &config.Config{
		AuthEnable:             true,
		AuthRegistrationEnable: false, // Disable normal registration initially
		AuthAccessTokenTTL:     60 * time.Minute,
		AuthRefreshTokenTTL:    24 * 30 * time.Hour,
	}

	// Initialize logger
	logger := slog.Default()

	// Initialize database connection - adapt to your actual DB initialization
	db, err := sql.Open("sqlite3", "plainq.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Initialize hasher and token manager
	hasher := hashkit.NewArgon2()
	tokenManager := jwtkit.New(jwtkit.WithHMACKey([]byte("your-secret-key"))) // Use proper key management

	// Initialize storage layers
	onboardingStorage := onboarding.NewSQLiteStorage(db)
	accountStorage := account.NewSQLiteStorage(db, hasher)
	rbacStorage := rbac.NewSQLiteStorage(db)

	// Initialize services
	onboardingService := onboarding.NewService(cfg, logger, hasher, tokenManager, onboardingStorage)
	accountService := account.NewService(cfg, logger, hasher, accountStorage)
	rbacService := rbac.NewService(cfg, logger, rbacStorage)

	// Check onboarding status at startup
	ctx := context.Background()
	needsOnboarding, err := onboardingService.NeedsOnboarding(ctx)
	if err != nil {
		log.Fatal("Failed to check onboarding status:", err)
	}

	if needsOnboarding {
		logger.Info("System needs onboarding - no admin users found")
		// Disable regular account registration until onboarding is complete
		cfg.AuthRegistrationEnable = false
	} else {
		logger.Info("System has been onboarded - admin users exist")
		// Enable regular account registration if desired
		cfg.AuthRegistrationEnable = true
	}

	// Setup router
	r := chi.NewRouter()

	// Add CORS and other basic middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.RedirectSlashes())

	// Add onboarding middleware to all routes except onboarding and health
	r.Use(middleware.RequireOnboarding(onboardingService))

	// Public routes (no authentication required)
	r.Route("/", func(r chi.Router) {
		// Onboarding endpoints (always available)
		r.Mount("/onboarding", onboardingService)

		// Health check (always available)
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		// Account endpoints (signin/signup)
		r.Mount("/account", accountService)
	})

	// Protected routes (require authentication)
	r.Route("/api", func(r chi.Router) {
		// Apply authentication middleware
		r.Use(middleware.AuthenticateJWT(tokenManager))

		// Admin-only routes
		r.Route("/admin", func(r chi.Router) {
			r.Use(middleware.RequireAdmin())
			r.Mount("/rbac", rbacService)
		})

		// Queue operations with permission checking
		r.Route("/queues/{queueID}", func(r chi.Router) {
			// Send message - requires send permission
			r.With(middleware.RequireQueueSendPermission(rbacService)).Post("/messages", sendMessageHandler)

			// Receive message - requires receive permission
			r.With(middleware.RequireQueueReceivePermission(rbacService)).Get("/messages", receiveMessageHandler)

			// Purge queue - requires purge permission or admin
			r.With(middleware.RequireAdminOrPermission(rbacService, middleware.PermissionPurge)).Delete("/messages", purgeQueueHandler)

			// Delete queue - requires delete permission or admin
			r.With(middleware.RequireAdminOrPermission(rbacService, middleware.PermissionDelete)).Delete("/", deleteQueueHandler)
		})

		// Role management (admin only)
		r.Route("/roles", func(r chi.Router) {
			r.Use(middleware.RequireAdmin())
			r.Get("/", listRolesHandler)
			r.Post("/", createRoleHandler)
		})
	})

	// Start server
	logger.Info("Starting PlainQ server", slog.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server failed:", err)
	}
}

// Example handlers - replace with your actual implementations
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for sending messages
	userInfo := middleware.MustGetUserFromContext(r.Context())
	logger := slog.Default()
	logger.Info("User sending message", slog.String("user_id", userInfo.UserID))
	w.WriteHeader(http.StatusOK)
}

func receiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for receiving messages
	userInfo := middleware.MustGetUserFromContext(r.Context())
	logger := slog.Default()
	logger.Info("User receiving message", slog.String("user_id", userInfo.UserID))
	w.WriteHeader(http.StatusOK)
}

func purgeQueueHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for purging queues
	userInfo := middleware.MustGetUserFromContext(r.Context())
	logger := slog.Default()
	logger.Info("User purging queue", slog.String("user_id", userInfo.UserID))
	w.WriteHeader(http.StatusOK)
}

func deleteQueueHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting queues
	userInfo := middleware.MustGetUserFromContext(r.Context())
	logger := slog.Default()
	logger.Info("User deleting queue", slog.String("user_id", userInfo.UserID))
	w.WriteHeader(http.StatusOK)
}

func listRolesHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for listing roles
	w.WriteHeader(http.StatusOK)
}

func createRoleHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating roles
	w.WriteHeader(http.StatusOK)
}