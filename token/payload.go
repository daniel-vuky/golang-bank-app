package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	NotBefore time.Time `json:"not_before"`
	ExpiredAt time.Time `json:"expired_at"`
	Issuer    string    `json:"issuer"`
	Subject   string    `json:"subject"`
	Audience  string    `json:"audience"`
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.ExpiredAt}, nil
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.IssuedAt}, nil
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: p.NotBefore}, nil
}

func (p Payload) GetIssuer() (string, error) {
	return p.Issuer, nil
}

func (p Payload) GetSubject() (string, error) {
	return p.Subject, nil
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{p.Audience}, nil
}

// NewPayload creates a new payload with specific username and duration
func NewPayload(username, role string, duration time.Duration) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        tokenID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		NotBefore: time.Now(),
		ExpiredAt: time.Now().Add(duration),
		Issuer:    "",
		Subject:   username,
		Audience:  "",
	}

	return payload, nil
}

// Valid checks if the token payload is valid or not
func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return errors.New("token expired")
	}
	return nil
}
