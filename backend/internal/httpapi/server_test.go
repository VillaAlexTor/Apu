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
		ID:          "session-id",
		CSRFHash:    csrfHash[:],
		MFAVerified: true,
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

func TestPlatformPermissionMatrix(t *testing.T) {
	platformAdmin := "platform_admin"
	supportAgent := "support_agent"
	tests := []struct {
		name     string
		session  identity.Session
		expected int
	}{
		{name: "platform admin with MFA", session: identity.Session{User: identity.User{PlatformRole: &platformAdmin}, MFAVerified: true}, expected: http.StatusNoContent},
		{name: "platform admin without MFA", session: identity.Session{User: identity.User{PlatformRole: &platformAdmin}}, expected: http.StatusForbidden},
		{name: "support agent", session: identity.Session{User: identity.User{PlatformRole: &supportAgent}, MFAVerified: true}, expected: http.StatusForbidden},
		{name: "tenant owner", session: identity.Session{Memberships: []identity.Membership{{Role: "tenant_owner"}}, MFAVerified: true}, expected: http.StatusForbidden},
		{name: "tenant admin", session: identity.Session{Memberships: []identity.Membership{{Role: "tenant_admin"}}, MFAVerified: true}, expected: http.StatusForbidden},
		{name: "guardian", session: identity.Session{Memberships: []identity.Membership{{Role: "guardian"}}}, expected: http.StatusForbidden},
	}
	server := &Server{platform: &stubPlatformStore{}}
	next := http.HandlerFunc(func(writer http.ResponseWriter, _ *http.Request) { writer.WriteHeader(http.StatusNoContent) })
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/v1/platform/summary", nil)
			request = request.WithContext(context.WithValue(request.Context(), sessionContextKey, test.session))
			response := httptest.NewRecorder()
			server.requirePlatformAdmin(next).ServeHTTP(response, request)
			if response.Code != test.expected {
				t.Fatalf("expected status %d, got %d", test.expected, response.Code)
			}
		})
	}
}

func TestPrivilegedTenantRolesRequireMFA(t *testing.T) {
	if !requiresMFA(identity.User{}, []identity.Membership{{Role: "tenant_owner"}}) {
		t.Fatal("tenant owners must require MFA")
	}
	if !requiresMFA(identity.User{}, []identity.Membership{{Role: "tenant_admin"}}) {
		t.Fatal("tenant administrators must require MFA")
	}
	if requiresMFA(identity.User{}, []identity.Membership{{Role: "guardian"}}) {
		t.Fatal("guardians must not require MFA")
	}
}

type stubPlatformStore struct{}

func (*stubPlatformStore) PlatformSummary(context.Context) (identity.PlatformSummary, error) {
	return identity.PlatformSummary{}, nil
}
func (*stubPlatformStore) PlatformTenants(context.Context) ([]identity.TenantRecord, error) {
	return nil, nil
}
func (*stubPlatformStore) CreateTenant(context.Context, identity.User, string, string) (identity.TenantRecord, error) {
	return identity.TenantRecord{}, nil
}
func (*stubPlatformStore) SetTenantStatus(context.Context, identity.User, string, string) error {
	return nil
}
func (*stubPlatformStore) PlatformUsers(context.Context) ([]identity.PlatformUser, error) {
	return nil, nil
}
func (*stubPlatformStore) CreateInvitation(context.Context, identity.User, string, string, string, string, []byte, time.Time) (identity.InvitationResult, error) {
	return identity.InvitationResult{}, nil
}
func (*stubPlatformStore) ActivateInvitation(context.Context, string, string) error { return nil }
func (*stubPlatformStore) SetUserStatus(context.Context, identity.User, string, string) error {
	return nil
}
func (*stubPlatformStore) RevokeUserSessions(context.Context, identity.User, string) (int64, error) {
	return 0, nil
}
func (*stubPlatformStore) SetMembershipRole(context.Context, identity.User, string, string, string) error {
	return nil
}
func (*stubPlatformStore) PlatformAudit(context.Context, int) ([]identity.AuditEvent, error) {
	return nil, nil
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
