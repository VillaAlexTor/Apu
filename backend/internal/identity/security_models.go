package identity

import "time"

type AuthChallenge struct {
	ID                     string
	Purpose                string
	User                   User
	PendingSecretEncrypted []byte
	MFASecretEncrypted     []byte
	AttemptCount           int
	ExpiresAt              time.Time
}
