package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/database"
	"github.com/VillaAlexTor/Apu/backend/internal/httpapi"
	"github.com/VillaAlexTor/Apu/backend/internal/identity"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "0.2.0"

type config struct {
	environment   string
	httpAddress   string
	databaseURL   string
	allowedOrigin string
	cookieSecure  bool
	sessionTTL    time.Duration
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("apu api stopped", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	rootContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := database.Open(rootContext, cfg.databaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	api, err := httpapi.New(identity.NewStore(pool), logger, httpapi.Config{
		AllowedOrigin: cfg.allowedOrigin,
		CookieSecure:  cfg.cookieSecure,
		SessionTTL:    cfg.sessionTTL,
	})
	if err != nil {
		return fmt.Errorf("create http api: %w", err)
	}

	server := &http.Server{
		Addr:              cfg.httpAddress,
		Handler:           api.Handler(healthHandler(), readinessHandler(pool)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("apu api listening", "address", cfg.httpAddress, "environment", cfg.environment)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case <-rootContext.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		logger.Info("apu api stopped gracefully")
		return nil
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("serve http: %w", err)
	}
}

func loadConfig() (config, error) {
	environment := envOrDefault("APU_ENV", "local")
	databaseURL := strings.TrimSpace(os.Getenv("APU_DATABASE_URL"))
	if databaseURL == "" {
		return config{}, errors.New("APU_DATABASE_URL is required")
	}
	sessionTTL, err := time.ParseDuration(envOrDefault("APU_SESSION_TTL", "12h"))
	if err != nil || sessionTTL < 15*time.Minute || sessionTTL > 7*24*time.Hour {
		return config{}, errors.New("APU_SESSION_TTL must be between 15m and 168h")
	}
	cookieSecure, err := strconv.ParseBool(envOrDefault("APU_COOKIE_SECURE", "false"))
	if err != nil {
		return config{}, errors.New("APU_COOKIE_SECURE must be true or false")
	}
	allowedOrigin := envOrDefault("APU_ALLOWED_ORIGIN", "http://localhost:3001")
	if environment == "production" && (!cookieSecure || !strings.HasPrefix(allowedOrigin, "https://")) {
		return config{}, errors.New("production requires secure cookies and an HTTPS allowed origin")
	}

	return config{
		environment:   environment,
		httpAddress:   envOrDefault("APU_HTTP_ADDR", ":8080"),
		databaseURL:   databaseURL,
		allowedOrigin: allowedOrigin,
		cookieSecure:  cookieSecure,
		sessionTTL:    sessionTTL,
	}, nil
}

func healthHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writeStatus(writer, http.StatusOK, "ok")
	}
}

func readinessHandler(pool *pgxpool.Pool) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(request.Context(), 2*time.Second)
		defer cancel()
		if err := pool.Ping(ctx); err != nil {
			writeStatus(writer, http.StatusServiceUnavailable, "database_unavailable")
			return
		}
		writeStatus(writer, http.StatusOK, "ready")
	}
}

func writeStatus(writer http.ResponseWriter, status int, state string) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(status)
	_, _ = fmt.Fprintf(writer, `{"service":"apu-api","status":%q,"version":%q}`+"\n", state, version)
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
