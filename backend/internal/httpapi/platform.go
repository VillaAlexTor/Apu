package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/identity"
	"github.com/VillaAlexTor/Apu/backend/internal/password"
	"github.com/go-chi/chi/v5"
)

var tenantSlugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

type platformStore interface {
	PlatformSummary(context.Context) (identity.PlatformSummary, error)
	PlatformTenants(context.Context) ([]identity.TenantRecord, error)
	CreateTenant(context.Context, identity.User, string, string) (identity.TenantRecord, error)
	SetTenantStatus(context.Context, identity.User, string, string) error
	PlatformUsers(context.Context) ([]identity.PlatformUser, error)
	CreateInvitation(context.Context, identity.User, string, string, string, string, []byte, time.Time) (identity.InvitationResult, error)
	ActivateInvitation(context.Context, string, string) error
	SetUserStatus(context.Context, identity.User, string, string) error
	RevokeUserSessions(context.Context, identity.User, string) (int64, error)
	SetMembershipRole(context.Context, identity.User, string, string, string) error
	PlatformAudit(context.Context, int) ([]identity.AuditEvent, error)
}

func platformStoreFrom(store identityStore) platformStore {
	platform, _ := store.(platformStore)
	return platform
}

func (server *Server) requirePlatformAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		session := sessionFromContext(request.Context())
		if session.User.PlatformRole == nil || *session.User.PlatformRole != "platform_admin" {
			writeError(writer, http.StatusForbidden, "forbidden", "No tienes permisos para operar la plataforma.")
			return
		}
		if server.platform == nil {
			server.internalError(writer, request, errors.New("platform store is unavailable"))
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func (server *Server) activateInvitation(writer http.ResponseWriter, request *http.Request) {
	if server.platform == nil {
		server.internalError(writer, request, errors.New("platform store is unavailable"))
		return
	}
	var input struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || len(input.Token) < 32 {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La invitación o la contraseña no son válidas.")
		return
	}
	passwordHash, err := password.Hash(input.Password)
	if err != nil {
		writeError(writer, http.StatusBadRequest, "weak_password", "La contraseña debe tener entre 12 y 1024 caracteres.")
		return
	}
	if err := server.platform.ActivateInvitation(request.Context(), input.Token, passwordHash); err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			writeError(writer, http.StatusBadRequest, "invitation_invalid", "La invitación no es válida o ya expiró.")
			return
		}
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": "activated"})
}

func (server *Server) platformSummary(writer http.ResponseWriter, request *http.Request) {
	summary, err := server.platform.PlatformSummary(request.Context())
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, summary)
}

func (server *Server) platformTenants(writer http.ResponseWriter, request *http.Request) {
	tenants, err := server.platform.PlatformTenants(request.Context())
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"items": tenants})
}

func (server *Server) createTenant(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La solicitud no tiene un formato válido.")
		return
	}
	input.Name = strings.TrimSpace(input.Name)
	input.Slug = strings.ToLower(strings.TrimSpace(input.Slug))
	if len(input.Name) == 0 || len(input.Name) > 160 || len(input.Slug) > 80 || !tenantSlugPattern.MatchString(input.Slug) {
		writeError(writer, http.StatusBadRequest, "invalid_tenant", "Revisa el nombre y el identificador de la organización.")
		return
	}
	tenant, err := server.platform.CreateTenant(request.Context(), sessionFromContext(request.Context()).User, input.Name, input.Slug)
	if err != nil {
		server.writePlatformError(writer, request, err, "Ya existe una organización con ese identificador.")
		return
	}
	writeJSON(writer, http.StatusCreated, tenant)
}

func (server *Server) setTenantStatus(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || !oneOf(input.Status, "active", "suspended", "deactivated") {
		writeError(writer, http.StatusBadRequest, "invalid_status", "El estado de la organización no es válido.")
		return
	}
	err := server.platform.SetTenantStatus(request.Context(), sessionFromContext(request.Context()).User, chi.URLParam(request, "tenantID"), input.Status)
	if err != nil {
		server.writePlatformError(writer, request, err, "No se pudo cambiar el estado de la organización.")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": input.Status})
}

func (server *Server) platformUsers(writer http.ResponseWriter, request *http.Request) {
	users, err := server.platform.PlatformUsers(request.Context())
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"items": users})
}

func (server *Server) createInvitation(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		TenantID    string `json:"tenantId"`
		Email       string `json:"email"`
		DisplayName string `json:"displayName"`
		Role        string `json:"role"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La solicitud no tiene un formato válido.")
		return
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.DisplayName = strings.TrimSpace(input.DisplayName)
	address, emailErr := mail.ParseAddress(input.Email)
	if emailErr != nil || strings.ToLower(address.Address) != input.Email || len(input.DisplayName) == 0 || len(input.DisplayName) > 160 || !oneOf(input.Role, "tenant_owner", "tenant_admin", "guardian") {
		writeError(writer, http.StatusBadRequest, "invalid_invitation", "Revisa el nombre, correo, organización y rol.")
		return
	}
	plainToken, tokenHash, err := identity.NewSecret()
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	result, err := server.platform.CreateInvitation(
		request.Context(),
		sessionFromContext(request.Context()).User,
		input.TenantID,
		input.Email,
		input.DisplayName,
		input.Role,
		tokenHash,
		server.now().UTC().Add(72*time.Hour),
	)
	if err != nil {
		server.writePlatformError(writer, request, err, "El usuario ya pertenece a esta organización o no puede ser invitado.")
		return
	}
	response := map[string]any{"user": result.User, "requiresActivation": result.RequiresActivation}
	if result.RequiresActivation {
		response["activationToken"] = plainToken
	}
	writeJSON(writer, http.StatusCreated, response)
}

func (server *Server) setUserStatus(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Status string `json:"status"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || !oneOf(input.Status, "active", "suspended", "deactivated") {
		writeError(writer, http.StatusBadRequest, "invalid_status", "El estado del usuario no es válido.")
		return
	}
	err := server.platform.SetUserStatus(request.Context(), sessionFromContext(request.Context()).User, chi.URLParam(request, "userID"), input.Status)
	if err != nil {
		server.writePlatformError(writer, request, err, "No se puede cambiar el estado de esta cuenta.")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"status": input.Status})
}

func (server *Server) revokeUserSessions(writer http.ResponseWriter, request *http.Request) {
	count, err := server.platform.RevokeUserSessions(request.Context(), sessionFromContext(request.Context()).User, chi.URLParam(request, "userID"))
	if err != nil {
		server.writePlatformError(writer, request, err, "No se pudieron revocar las sesiones.")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]int64{"revokedSessions": count})
}

func (server *Server) setMembershipRole(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Role string `json:"role"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || !oneOf(input.Role, "tenant_owner", "tenant_admin", "guardian") {
		writeError(writer, http.StatusBadRequest, "invalid_role", "El rol no es válido.")
		return
	}
	err := server.platform.SetMembershipRole(
		request.Context(),
		sessionFromContext(request.Context()).User,
		chi.URLParam(request, "tenantID"),
		chi.URLParam(request, "userID"),
		input.Role,
	)
	if err != nil {
		server.writePlatformError(writer, request, err, "No se pudo actualizar el rol.")
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{"role": input.Role})
}

func (server *Server) platformAudit(writer http.ResponseWriter, request *http.Request) {
	events, err := server.platform.PlatformAudit(request.Context(), 30)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]any{"items": events})
}

func (server *Server) writePlatformError(writer http.ResponseWriter, request *http.Request, err error, conflictMessage string) {
	switch {
	case errors.Is(err, identity.ErrNotFound):
		writeError(writer, http.StatusNotFound, "not_found", "El recurso solicitado no existe.")
	case errors.Is(err, identity.ErrConflict):
		writeError(writer, http.StatusConflict, "conflict", conflictMessage)
	default:
		server.internalError(writer, request, err)
	}
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
