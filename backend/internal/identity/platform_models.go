package identity

import (
	"encoding/json"
	"time"
)

type PlatformSummary struct {
	TenantsTotal       int `json:"tenantsTotal"`
	TenantsActive      int `json:"tenantsActive"`
	UsersActive        int `json:"usersActive"`
	UsersInvited       int `json:"usersInvited"`
	UsersSuspended     int `json:"usersSuspended"`
	PendingInvitations int `json:"pendingInvitations"`
}

type TenantRecord struct {
	ID             string    `json:"id"`
	Slug           string    `json:"slug"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	OwnerName      *string   `json:"ownerName,omitempty"`
	OwnerEmail     *string   `json:"ownerEmail,omitempty"`
	UserCount      int       `json:"userCount"`
	ActiveSessions int       `json:"activeSessions"`
	CreatedAt      time.Time `json:"createdAt"`
}

type PlatformUser struct {
	ID             string       `json:"id"`
	Email          string       `json:"email"`
	DisplayName    string       `json:"displayName"`
	Status         string       `json:"status"`
	PlatformRole   *string      `json:"platformRole,omitempty"`
	Memberships    []Membership `json:"memberships"`
	ActiveSessions int          `json:"activeSessions"`
	CreatedAt      time.Time    `json:"createdAt"`
}

type AuditEvent struct {
	ID           int64           `json:"id"`
	ActorEmail   *string         `json:"actorEmail,omitempty"`
	TenantName   *string         `json:"tenantName,omitempty"`
	Action       string          `json:"action"`
	ResourceType string          `json:"resourceType"`
	ResourceID   *string         `json:"resourceId,omitempty"`
	Metadata     json.RawMessage `json:"metadata"`
	OccurredAt   time.Time       `json:"occurredAt"`
}

type InvitationResult struct {
	User               PlatformUser `json:"user"`
	RequiresActivation bool         `json:"requiresActivation"`
}
