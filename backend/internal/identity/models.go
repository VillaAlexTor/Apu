package identity

import "time"

type User struct {
	ID           string
	Email        string
	DisplayName  string
	PasswordHash string
	Status       string
	PlatformRole *string
}

type Membership struct {
	TenantID   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
	Role       string `json:"role"`
}

type Session struct {
	ID          string
	User        User
	CSRFHash    []byte
	ExpiresAt   time.Time
	Memberships []Membership
}

type PublicUser struct {
	ID           string       `json:"id"`
	Email        string       `json:"email"`
	DisplayName  string       `json:"displayName"`
	PlatformRole *string      `json:"platformRole,omitempty"`
	Memberships  []Membership `json:"memberships"`
}

func (session Session) PublicUser() PublicUser {
	return PublicUser{
		ID:           session.User.ID,
		Email:        session.User.Email,
		DisplayName:  session.User.DisplayName,
		PlatformRole: session.User.PlatformRole,
		Memberships:  session.Memberships,
	}
}

func Destination(user User, memberships []Membership) string {
	if user.PlatformRole != nil {
		return "/platform"
	}
	for _, membership := range memberships {
		if membership.Role == "tenant_owner" || membership.Role == "tenant_admin" {
			return "/admin"
		}
	}
	for _, membership := range memberships {
		if membership.Role == "guardian" {
			return "/familia"
		}
	}
	return "/acceso-pendiente"
}
