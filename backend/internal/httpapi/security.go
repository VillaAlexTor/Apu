package httpapi

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/VillaAlexTor/Apu/backend/internal/identity"
	"github.com/VillaAlexTor/Apu/backend/internal/mfa"
	"github.com/VillaAlexTor/Apu/backend/internal/password"
)

type ResetMailer interface {
	SendPasswordReset(context.Context, string, string) error
}

type securityStore interface {
	CreateAuthChallenge(context.Context, string, string, []byte, time.Time) error
	AuthChallengeByToken(context.Context, string, string) (identity.AuthChallenge, error)
	SetAuthChallengeSecret(context.Context, string, []byte) error
	RecordAuthChallengeFailure(context.Context, string) error
	ConsumeAuthChallenge(context.Context, string) error
	EnableMFA(context.Context, string, []byte, int64, [][]byte) error
	ConsumeTOTP(context.Context, string, int64) error
	ConsumeRecoveryCode(context.Context, string, []byte) error
	CreateMFASession(context.Context, identity.User, []byte, []byte, time.Time) (identity.Session, error)
	CreatePasswordReset(context.Context, string, []byte, time.Time) error
	ResetPassword(context.Context, string, string) error
	ChangePassword(context.Context, string, string) error
}

func securityStoreFrom(store identityStore) securityStore {
	security, _ := store.(securityStore)
	return security
}

func requiresMFA(user identity.User, memberships []identity.Membership) bool {
	if user.PlatformRole != nil && *user.PlatformRole == "platform_admin" {
		return true
	}
	for _, membership := range memberships {
		if membership.Role == "tenant_owner" || membership.Role == "tenant_admin" {
			return true
		}
	}
	return false
}

func (server *Server) createSessionResponse(
	writer http.ResponseWriter,
	request *http.Request,
	user identity.User,
	memberships []identity.Membership,
	mfaVerified bool,
	extra map[string]any,
) {
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
	var session identity.Session
	if mfaVerified {
		session, err = server.security.CreateMFASession(request.Context(), user, tokenHash, csrfHash, expiresAt)
	} else {
		session, err = server.store.CreateSession(request.Context(), user, tokenHash, csrfHash, expiresAt)
	}
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	session.Memberships = memberships
	session.MFAVerified = mfaVerified
	server.limiter.reset(remoteIP(request.RemoteAddr))
	server.setAuthCookies(writer, sessionToken, csrfToken, expiresAt)
	body := map[string]any{
		"user":      session.PublicUser(),
		"redirect":  identity.Destination(user, memberships),
		"csrfToken": csrfToken,
	}
	for key, value := range extra {
		body[key] = value
	}
	writeJSON(writer, http.StatusOK, body)
}

func (server *Server) startMFAEnrollment(writer http.ResponseWriter, request *http.Request) {
	if server.security == nil {
		server.internalError(writer, request, errors.New("security store is unavailable"))
		return
	}
	var input struct {
		Challenge string `json:"challenge"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || len(input.Challenge) < 32 {
		server.invalidMFAChallenge(writer)
		return
	}
	challenge, err := server.security.AuthChallengeByToken(request.Context(), input.Challenge, "mfa_enroll")
	if err != nil || challenge.User.MFAEnabled {
		server.invalidMFAChallenge(writer)
		return
	}
	var secret string
	if len(challenge.PendingSecretEncrypted) == 0 {
		secret, err = mfa.GenerateSecret()
		if err == nil {
			var encrypted []byte
			encrypted, err = server.mfaEncryptor.Encrypt(secret)
			if err == nil {
				err = server.security.SetAuthChallengeSecret(request.Context(), challenge.ID, encrypted)
			}
		}
	} else {
		secret, err = server.mfaEncryptor.Decrypt(challenge.PendingSecretEncrypted)
	}
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	writeJSON(writer, http.StatusOK, map[string]string{
		"secret":          secret,
		"provisioningUri": mfa.ProvisioningURI("APU", challenge.User.Email, secret),
	})
}

func (server *Server) confirmMFAEnrollment(writer http.ResponseWriter, request *http.Request) {
	if server.security == nil {
		server.internalError(writer, request, errors.New("security store is unavailable"))
		return
	}
	var input struct {
		Challenge string `json:"challenge"`
		Code      string `json:"code"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		server.invalidMFAChallenge(writer)
		return
	}
	challenge, err := server.security.AuthChallengeByToken(request.Context(), input.Challenge, "mfa_enroll")
	if err != nil || len(challenge.PendingSecretEncrypted) == 0 || challenge.User.MFAEnabled {
		server.invalidMFAChallenge(writer)
		return
	}
	secret, err := server.mfaEncryptor.Decrypt(challenge.PendingSecretEncrypted)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	step, valid := mfa.VerifyStep(secret, input.Code, server.now())
	if !valid {
		_ = server.security.RecordAuthChallengeFailure(request.Context(), challenge.ID)
		writeError(writer, http.StatusUnauthorized, "mfa_code_invalid", "El código no es válido o ya expiró.")
		return
	}
	recoveryCodes, recoveryHashes, err := mfa.GenerateRecoveryCodes()
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	if err := server.security.EnableMFA(request.Context(), challenge.User.ID, challenge.PendingSecretEncrypted, step, recoveryHashes); err != nil {
		server.writeSecurityError(writer, request, err)
		return
	}
	if err := server.security.ConsumeAuthChallenge(request.Context(), challenge.ID); err != nil {
		server.writeSecurityError(writer, request, err)
		return
	}
	memberships, err := server.store.Memberships(request.Context(), challenge.User.ID)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	challenge.User.MFAEnabled = true
	server.createSessionResponse(writer, request, challenge.User, memberships, true, map[string]any{"recoveryCodes": recoveryCodes})
}

func (server *Server) verifyMFA(writer http.ResponseWriter, request *http.Request) {
	if server.security == nil {
		server.internalError(writer, request, errors.New("security store is unavailable"))
		return
	}
	var input struct {
		Challenge string `json:"challenge"`
		Code      string `json:"code"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		server.invalidMFAChallenge(writer)
		return
	}
	challenge, err := server.security.AuthChallengeByToken(request.Context(), input.Challenge, "mfa_verify")
	if err != nil || !challenge.User.MFAEnabled || len(challenge.MFASecretEncrypted) == 0 {
		server.invalidMFAChallenge(writer)
		return
	}
	secret, err := server.mfaEncryptor.Decrypt(challenge.MFASecretEncrypted)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	step, validTOTP := mfa.VerifyStep(secret, input.Code, server.now())
	verified := false
	if validTOTP {
		verified = server.security.ConsumeTOTP(request.Context(), challenge.User.ID, step) == nil
	} else if strings.TrimSpace(input.Code) != "" {
		verified = server.security.ConsumeRecoveryCode(request.Context(), challenge.User.ID, mfa.HashRecoveryCode(input.Code)) == nil
	}
	if !verified {
		_ = server.security.RecordAuthChallengeFailure(request.Context(), challenge.ID)
		writeError(writer, http.StatusUnauthorized, "mfa_code_invalid", "El código no es válido, ya fue utilizado o expiró.")
		return
	}
	if err := server.security.ConsumeAuthChallenge(request.Context(), challenge.ID); err != nil {
		server.writeSecurityError(writer, request, err)
		return
	}
	memberships, err := server.store.Memberships(request.Context(), challenge.User.ID)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	server.createSessionResponse(writer, request, challenge.User, memberships, true, nil)
}

func (server *Server) requestPasswordReset(writer http.ResponseWriter, request *http.Request) {
	startedAt := time.Now()
	defer func() {
		if remaining := 350*time.Millisecond - time.Since(startedAt); remaining > 0 {
			time.Sleep(remaining)
		}
	}()
	var input struct {
		Email string `json:"email"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La solicitud no tiene un formato válido.")
		return
	}
	response := map[string]any{"status": "accepted"}
	clientKey := remoteIP(request.RemoteAddr)
	if !server.resetLimiter.allow(clientKey) || server.security == nil {
		writeJSON(writer, http.StatusAccepted, response)
		return
	}
	email := strings.ToLower(strings.TrimSpace(input.Email))
	user, err := server.store.FindUserByEmail(request.Context(), email)
	if err == nil && user.Status == "active" && user.PasswordHash != "" {
		token, tokenHash, tokenErr := identity.NewSecret()
		if tokenErr != nil {
			server.internalError(writer, request, tokenErr)
			return
		}
		if err := server.security.CreatePasswordReset(request.Context(), user.ID, tokenHash, server.now().UTC().Add(30*time.Minute)); err != nil {
			server.internalError(writer, request, err)
			return
		}
		resetURL := strings.TrimRight(server.config.PasswordResetBaseURL, "/") + "/restablecer?token=" + url.QueryEscape(token)
		if server.config.ResetMailer != nil {
			go func(recipient, link string) {
				ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
				defer cancel()
				if err := server.config.ResetMailer.SendPasswordReset(ctx, recipient, link); err != nil {
					server.logger.Error("password reset delivery failed", "error", err)
				}
			}(user.Email, resetURL)
		}
		if server.config.ExposeDevelopmentTokens {
			response["developmentResetToken"] = token
		}
	}
	writeJSON(writer, http.StatusAccepted, response)
}

func (server *Server) confirmPasswordReset(writer http.ResponseWriter, request *http.Request) {
	if server.security == nil {
		server.internalError(writer, request, errors.New("security store is unavailable"))
		return
	}
	var input struct {
		Token       string `json:"token"`
		NewPassword string `json:"newPassword"`
	}
	if err := decodeJSON(writer, request, &input); err != nil || len(input.Token) < 32 {
		writeError(writer, http.StatusBadRequest, "reset_invalid", "El enlace no es válido o ya expiró.")
		return
	}
	passwordHash, err := password.Hash(input.NewPassword)
	if err != nil {
		server.writePasswordPolicyError(writer, err)
		return
	}
	if err := server.security.ResetPassword(request.Context(), input.Token, passwordHash); err != nil {
		if errors.Is(err, identity.ErrNotFound) {
			writeError(writer, http.StatusBadRequest, "reset_invalid", "El enlace no es válido o ya expiró.")
			return
		}
		server.internalError(writer, request, err)
		return
	}
	server.clearAuthCookies(writer)
	writeJSON(writer, http.StatusOK, map[string]string{"status": "password_reset"})
}

func (server *Server) changePassword(writer http.ResponseWriter, request *http.Request) {
	if server.security == nil {
		server.internalError(writer, request, errors.New("security store is unavailable"))
		return
	}
	var input struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}
	if err := decodeJSON(writer, request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, "invalid_request", "La solicitud no tiene un formato válido.")
		return
	}
	session := sessionFromContext(request.Context())
	valid, err := password.Verify(input.CurrentPassword, session.User.PasswordHash)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	if !valid {
		writeError(writer, http.StatusUnauthorized, "current_password_invalid", "La contraseña actual no es válida.")
		return
	}
	same, err := password.Verify(input.NewPassword, session.User.PasswordHash)
	if err != nil {
		server.internalError(writer, request, err)
		return
	}
	if same {
		writeError(writer, http.StatusBadRequest, "password_reused", "La nueva contraseña debe ser diferente de la actual.")
		return
	}
	passwordHash, err := password.Hash(input.NewPassword)
	if err != nil {
		server.writePasswordPolicyError(writer, err)
		return
	}
	if err := server.security.ChangePassword(request.Context(), session.User.ID, passwordHash); err != nil {
		server.internalError(writer, request, err)
		return
	}
	server.clearAuthCookies(writer)
	writer.WriteHeader(http.StatusNoContent)
}

func (server *Server) invalidMFAChallenge(writer http.ResponseWriter) {
	writeError(writer, http.StatusUnauthorized, "mfa_challenge_invalid", "La verificación expiró. Inicia sesión nuevamente.")
}

func (server *Server) writePasswordPolicyError(writer http.ResponseWriter, err error) {
	if errors.Is(err, password.ErrCommonPassword) {
		writeError(writer, http.StatusBadRequest, "password_common", "Elige una contraseña menos común.")
		return
	}
	writeError(writer, http.StatusBadRequest, "password_weak", "Usa al menos 15 caracteres y evita contraseñas conocidas.")
}

func (server *Server) writeSecurityError(writer http.ResponseWriter, request *http.Request, err error) {
	if errors.Is(err, identity.ErrNotFound) || errors.Is(err, identity.ErrConflict) {
		server.invalidMFAChallenge(writer)
		return
	}
	server.internalError(writer, request, err)
}
