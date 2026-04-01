package interfaces

import (
	"context"
	"time"
)

type SessionRepository interface {
	Save(ctx context.Context, accountID, refToken string, expiresAt time.Time) error
	Delete(ctx context.Context, refToken string) error
}
