package identity

import "testing"

func TestDestination(t *testing.T) {
	platformRole := "platform_admin"
	tests := []struct {
		name        string
		user        User
		memberships []Membership
		expected    string
	}{
		{name: "platform", user: User{PlatformRole: &platformRole}, expected: "/platform"},
		{name: "tenant admin", memberships: []Membership{{Role: "tenant_admin"}}, expected: "/admin"},
		{name: "guardian", memberships: []Membership{{Role: "guardian"}}, expected: "/familia"},
		{name: "pending", expected: "/acceso-pendiente"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if actual := Destination(test.user, test.memberships); actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}
