package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const version = "0.1.0"

var databaseDial = func(network, address string, timeout time.Duration) (io.Closer, error) {
	return net.DialTimeout(network, address, timeout)
}

type config struct {
	environment string
	httpAddress string
	databaseURL string
}

type response struct {
	Service     string `json:"service"`
	Status      string `json:"status"`
	Version     string `json:"version,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func main() {
	cfg := loadConfig()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", writeJSON(http.StatusOK, response{
		Service: "apu-api",
		Status:  "ok",
		Version: version,
	}))
	mux.HandleFunc("GET /readyz", readinessHandler(cfg))
	mux.HandleFunc("GET /api/v1", writeJSON(http.StatusOK, response{
		Service:     "apu-api",
		Status:      "bootstrap",
		Version:     version,
		Environment: cfg.environment,
	}))

	server := &http.Server{
		Addr:              cfg.httpAddress,
		Handler:           requestLogger(logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("apu api listening", "address", cfg.httpAddress, "environment", cfg.environment)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("api server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-shutdownSignal.Done()

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	logger.Info("apu api stopped")
}

func loadConfig() config {
	return config{
		environment: envOrDefault("APU_ENV", "local"),
		httpAddress: envOrDefault("APU_HTTP_ADDR", ":8080"),
		databaseURL: strings.TrimSpace(os.Getenv("APU_DATABASE_URL")),
	}
}

func readinessHandler(cfg config) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		if cfg.databaseURL == "" {
			writeJSON(http.StatusServiceUnavailable, response{
				Service: "apu-api",
				Status:  "database_not_configured",
			})(writer, nil)
			return
		}
		if err := databaseReachable(cfg.databaseURL); err != nil {
			writeJSON(http.StatusServiceUnavailable, response{
				Service: "apu-api",
				Status:  "database_unavailable",
			})(writer, nil)
			return
		}

		writeJSON(http.StatusOK, response{
			Service: "apu-api",
			Status:  "ready",
			Version: version,
		})(writer, nil)
	}
}

func databaseReachable(databaseURL string) error {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return err
	}

	hostname := parsedURL.Hostname()
	if hostname == "" {
		return errors.New("database hostname is missing")
	}
	port := parsedURL.Port()
	if port == "" {
		port = "5432"
	}

	connection, err := databaseDial("tcp", net.JoinHostPort(hostname, port), 2*time.Second)
	if err != nil {
		return err
	}
	return connection.Close()
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		next.ServeHTTP(writer, request)
		logger.Info("http request",
			"method", request.Method,
			"path", request.URL.Path,
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	})
}

func writeJSON(status int, body response) http.HandlerFunc {
	return func(writer http.ResponseWriter, _ *http.Request) {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		writer.Header().Set("Cache-Control", "no-store")
		writer.WriteHeader(status)
		if err := json.NewEncoder(writer).Encode(body); err != nil {
			slog.Error("could not encode response", "error", err)
		}
	}
}

func envOrDefault(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}
