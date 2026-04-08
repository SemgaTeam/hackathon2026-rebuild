package entities

import (
	"time"

	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	RefToken  string
	ExpiresAt time.Time
	Revoked   bool
}

func NewSession(id, accountId uuid.UUID, refToken string, expiresAt time.Time, revoked bool) (*Session, error) {
	if refToken == "" {
		return nil, domainErrors.NewError("incorrect refToken: refresh token cannot be empty")
	}

	if time.Now().After(expiresAt) {
		return nil, domainErrors.NewError("invalid expiresAt: expiresAt must be in the future")
	}

	return &Session{
		ID:        id,
		AccountID: accountId,
		RefToken:  refToken,
		ExpiresAt: expiresAt,
		Revoked:   revoked,
	}, nil

}
