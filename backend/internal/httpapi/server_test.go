package httpapi

import (
	"context"
	"crypto/sha256"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/identity"
)

type stubIdentityStore struct {
	session identity.Session
	revoked bool
	user    *identity.User
}

func (store *stubIdentityStore) FindUserByEmail(context.Context, string) (identity.User, error) {
	if store.user != nil {
		return *store.user, nil
	}
	return identity.User{}, identity.ErrNotFound
}

func (store *stubIdentityStore) Memberships(context.Context, string) ([]identity.Membership, error) {
	return nil, nil
}

func (store *stubIdentityStore) CreateSession(context.Context, identity.User, []byte, []byte, time.Time) (identity.Session, error) {
	return identity.Session{}, nil
}

func (store *stubIdentityStore) SessionByToken(_ context.Context, token string) (identity.Session, error) {
	if token != "valid-session" {
		return identity.Session{}, identity.ErrNotFound
	}
	return store.session, nil
}

func (store *stubIdentityStore) RevokeSession(context.Context, identity.Session) error {
	store.revoked = true
	return nil
}

func TestLogoutRejectsInvalidCSRF(t *testing.T) {
	csrfHash := sha256.Sum256([]byte("correct-csrf-token"))
	platformRole := "platform_admin"
	store := &stubIdentityStore{session: identity.Session{
		ID:       "session-id",
		CSRFHash: csrfHash[:],
		User: identity.User{
			ID:           "user-id",
			Email:        "admin@apu.local",
			DisplayName:  "Admin",
			Status:       "active",
			PlatformRole: &platformRole,
		},
	}}
	handler := testHandler(t, store)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: "apu_session", Value: "valid-session"})
	request.Header.Set("Origin", "http://localhost:3001")
	request.Header.Set("X-CSRF-Token", "incorrect-csrf-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
	if store.revoked {
		t.Fatal("session must not be revoked when the CSRF token is invalid")
	}
}

func TestLogoutAcceptsValidCSRF(t *testing.T) {
	csrfHash := sha256.Sum256([]byte("correct-csrf-token"))
	store := &stubIdentityStore{session: identity.Session{
		ID:       "session-id",
		CSRFHash: csrfHash[:],
		User:     identity.User{ID: "user-id", Status: "active"},
	}}
	handler := testHandler(t, store)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
	request.AddCookie(&http.Cookie{Name: "apu_session", Value: "valid-session"})
	request.Header.Set("Origin", "http://localhost:3001")
	request.Header.Set("X-CSRF-Token", "correct-csrf-token")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, response.Code)
	}
	if !store.revoked {
		t.Fatal("expected the session to be revoked")
	}
}

func TestCORSRejectsUnknownOrigin(t *testing.T) {
	handler := testHandler(t, &stubIdentityStore{})
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", nil)
	request.Header.Set("Origin", "https://malicioso.invalid")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
}

func TestInvitedUserLoginReturnsInvalidCredentials(t *testing.T) {
	store := &stubIdentityStore{user: &identity.User{
		ID:           "invited-user",
		Email:        "invited@apu.local",
		DisplayName:  "Invited",
		PasswordHash: "",
		Status:       "invited",
	}}
	handler := testHandler(t, store)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(`{"email":"invited@apu.local","password":"not-active-yet"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", "http://localhost:3001")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, response.Code)
	}
}

func TestPlatformRoutesRejectNonAdministrator(t *testing.T) {
	csrfHash := sha256.Sum256([]byte("csrf-token"))
	store := &stubIdentityStore{session: identity.Session{
		ID:       "session-id",
		CSRFHash: csrfHash[:],
		User: identity.User{
			ID:          "guardian-id",
			Email:       "guardian@apu.local",
			DisplayName: "Guardian",
			Status:      "active",
		},
	}}
	handler := testHandler(t, store)
	request := httptest.NewRequest(http.MethodGet, "/api/v1/platform/summary", nil)
	request.AddCookie(&http.Cookie{Name: "apu_session", Value: "valid-session"})
	request.Header.Set("Origin", "http://localhost:3001")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, response.Code)
	}
}

func testHandler(t *testing.T, store identityStore) http.Handler {
	t.Helper()
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	server, err := New(store, logger, Config{
		AllowedOrigin: "http://localhost:3001",
		SessionTTL:    12 * time.Hour,
	})
	if err != nil {
		t.Fatalf("create server: %v", err)
	}
	noContent := func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	}
	return server.Handler(noContent, noContent)
}
