package token

import "time"

// Maker interface for making tokens
type Maker interface {
	// CreateToken create a new token for a specific username and duration
	CreateToken(username, role string, duration time.Duration) (string, *Payload, error)

	// VerifyToken checks if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
