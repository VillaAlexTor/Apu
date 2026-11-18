package httpapi

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/identity"
	"github.com/VillaAlexTor/Apu/backend/internal/password"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const maxRequestBody = 1 << 20

type identityStore interface {
	FindUserByEmail(context.Context, string) (identity.User, error)
	Memberships(context.Context, string) ([]identity.Membership, error)
	CreateSession(context.Context, identity.User, []byte, []byte, time.Time) (identity.Session, error)
	SessionByToken(context.Context, string) (identity.Session, error)
	RevokeSession(context.Context, identity.Session) error
}

type Config struct {
	AllowedOrigin string
	CookieSecure  bool
	SessionTTL    time.Duration
}

type Server struct {
	store             identityStore
	logger            *slog.Logger
	config            Config
	limiter           *loginLimiter
	dummyPasswordHash string
	now               func() time.Time
}

type contextKey string

const sessionContextKey contextKey = "session"

func New(store identityStore, logger *slog.Logger, config Config) (*Server, error) {
	dummyHash, err := password.Hash("invalid-account-password")
	if err != nil {
		return nil, err
	}
	return &Server{
		store:             store,
		logger:            logger,
		config:            config,
		limiter:           newLoginLimiter(8, 15*time.Minute),
		dummyPasswordHash: dummyHash,
		now:               time.Now,
	}, nil
}

func (server *Server) Handler(healthHandler, readinessHandler http.HandlerFunc) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(server.cors)
	router.Use(server.requestLogger)

	router.Get("/healthz", healthHandler)
	router.Get("/readyz", readinessHandler)
	router.Get("/api/v1", func(writer http.ResponseWriter, request *http.Request) {
		writeJSON(writer, http.StatusOK, map[string]string{
			"service": "apu-api",
			"status":  "ready",
			"version": "0.2.0",
		})
	})
	router.Route("/api/v1/auth", func(router chi.Router) {
		router.Post("/login", server.login)
		router.Group(func(router chi.Router) {
			router.Use(server.authenticate)
			router.Get("/me", server.me)
			router.With(server.requireCSRF).Post("/logout", server.logout)
		})
	})
	return router
}

func (server *Server) login(writer http.ResponseWriter, request *http.Request) {
	clientKey := remoteIP(request.RemoteAddr)
	if !server.limiter.allow(clientKey) {
		writer.Header().Set("Retry-After", "900")
		writeError(writer, http.StatusTooManyRequests, "rate_limited", "Demasiados intentos. Intenta nuevamente más tarde.")
		return
	}

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La solicitud no tiene un formato válido.")
		return
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if email == "" || len(email) > 320 || len(input.Password) == 0 || len(input.Password) > 1024 {
		server.invalidCredentials(writer)
		return
	}

	user, err := server.store.FindUserByEmail(request.Context(), email)
	passwordHash := server.dummyPasswordHash
	if err == nil {
		passwordHash = user.PasswordHash
	} else if !errors.Is(err, identity.ErrNotFound) {
		server.internalError(writer, request, err)
		return
	}
	validPassword, verifyErr := password.Verify(input.Password, passwordHash)
	if verifyErr != nil {
		server.internalError(writer, request, verifyErr)
		return
	}
	if err != nil || !validPassword || user.Status != "active" {
		server.invalidCredentials(writer)
		return
	}

	memberships, err := server.store.Memberships(request.Context(), user.ID)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	sessionToken, tokenHash, err := identity.NewSecret()
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	csrfToken, csrfHash, err := identity.NewSecret()
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	expiresAt := server.now().UTC().Add(server.config.SessionTTL)
	session, err := server.store.CreateSession(request.Context(), user, tokenHash, csrfHash, expiresAt)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	session.Memberships = memberships
	server.setAuthCookies(writer, sessionToken, csrfToken, expiresAt)
	server.limiter.reset(clientKey)

	writeJSON(writer, http.StatusOK, map[string]any{
		"user":      session.PublicUser(),
		"redirect":  identity.Destination(user, memberships),
		"csrfToken": csrfToken,
	})
}

func (server *Server) me(writer http.ResponseWriter, request *http.Request) {
	session := sessionFromContext(request.Context())
	writeJSON(writer, http.StatusOK, map[string]any{
		"user":     session.PublicUser(),
		"redirect": identity.Destination(session.User, session.Memberships),
	})
}

func (server *Server) logout(writer http.ResponseWriter, request *http.Request) {
	session := sessionFromContext(request.Context())
	if err := server.store.RevokeSession(request.Context(), session); err != nil {
		server.internalError(writer, request, err)
		return
	}
	server.clearAuthCookies(writer)
	writer.WriteHeader(http.StatusNoContent)
}

func (server *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		cookie, err := request.Cookie(server.sessionCookieName())
		if err != nil || cookie.Value == "" {
			writeError(writer, http.StatusUnauthorized, "unauthorized", "La sesión no es válida o ha expirado.")
			return
		}
		session, err := server.store.SessionByToken(request.Context(), cookie.Value)
		if err != nil {
			if !errors.Is(err, identity.ErrNotFound) {
				server.internalError(writer, request, err)
				return
			}
			server.clearAuthCookies(writer)
			writeError(writer, http.StatusUnauthorized, "unauthorized", "La sesión no es válida o ha expirado.")
			return
		}
		ctx := context.WithValue(request.Context(), sessionContextKey, session)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func (server *Server) requireCSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session := sessionFromContext(request.Context())
		provided := request.Header.Get("X-CSRF-Token")
		digest := sha256.Sum256([]byte(provided))
		if provided == "" || subtle.ConstantTimeCompare(digest[:], session.CSRFHash) != 1 {
			writeError(writer, http.StatusForbidden, "csrf_invalid", "La verificación de seguridad no es válida.")
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (server *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		if origin != "" && origin != server.config.AllowedOrigin {
			writeError(writer, http.StatusForbidden, "origin_forbidden", "El origen de la solicitud no está permitido.")
			return
		}
		if origin == server.config.AllowedOrigin {
			writer.Header().Set("Access-Control-Allow-Origin", origin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			writer.Header().Add("Vary", "Origin")
		}
		if request.Method == http.MethodOptions {
			writer.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (server *Server) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		startedAt := time.Now()
		next.ServeHTTP(writer, request)
		server.logger.Info("http request",
			"request_id", middleware.GetReqID(request.Context()),
			"method", request.Method,
			"path", request.URL.Path,
			"duration_ms", time.Since(startedAt).Milliseconds(),
		)
	})
}

func (server *Server) setAuthCookies(writer http.ResponseWriter, sessionToken, csrfToken string, expiresAt time.Time) {
	http.SetCookie(writer, &http.Cookie{
		Name:     server.sessionCookieName(),
		Value:    sessionToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   server.config.CookieSecure,
		SameSite: http.SameSiteLaxMode,
	})
	http.SetCookie(writer, &http.Cookie{
		Name:     server.csrfCookieName(),
		Value:    csrfToken,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: false,
		Secure:   server.config.CookieSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func (server *Server) clearAuthCookies(writer http.ResponseWriter) {
	for _, cookieName := range []string{server.sessionCookieName(), server.csrfCookieName()} {
		http.SetCookie(writer, &http.Cookie{
			Name:     cookieName,
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Expires:  time.Unix(1, 0),
			HttpOnly: cookieName == server.sessionCookieName(),
			Secure:   server.config.CookieSecure,
			SameSite: http.SameSiteLaxMode,
		})
	}
}

func (server *Server) sessionCookieName() string {
	if server.config.CookieSecure {
		return "__Host-apu_session"
	}
	return "apu_session"
}

func (server *Server) csrfCookieName() string {
	if server.config.CookieSecure {
		return "__Host-apu_csrf"
	}
	return "apu_csrf"
}

func (server *Server) invalidCredentials(writer http.ResponseWriter) {
	writeError(writer, http.StatusUnauthorized, "invalid_credentials", "El correo o la contraseña no son válidos.")
}

func (server *Server) internalError(writer http.ResponseWriter, request *http.Request, err error) {
	server.logger.Error("request failed", "error", err, "request_id", middleware.GetReqID(request.Context()))
	writeError(writer, http.StatusInternalServerError, "internal_error", "No pudimos completar la solicitud.")
}

func decodeJSON(writer http.ResponseWriter, request *http.Request, destination any) error {
	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBody)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain exactly one JSON value")
	}
	return nil
}

func writeJSON(writer http.ResponseWriter, status int, body any) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.Header().Set("Cache-Control", "no-store")
	writer.WriteHeader(status)
	if err := json.NewEncoder(writer).Encode(body); err != nil {
		slog.Error("encode response", "error", err)
	}
}

func writeError(writer http.ResponseWriter, status int, code, message string) {
	writeJSON(writer, status, map[string]any{
		"error": map[string]string{"code": code, "message": message},
	})
}

func sessionFromContext(ctx context.Context) identity.Session {
	session, _ := ctx.Value(sessionContextKey).(identity.Session)
	return session
}

func remoteIP(remoteAddress string) string {
	host, _, err := net.SplitHostPort(remoteAddress)
	if err != nil {
		return remoteAddress
	}
	return host
}
