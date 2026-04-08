package interfaces

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type SessionRepository interface {
	FindByToken(ctx context.Context, refTokens string) (*entities.Session, error)
	Save(ctx context.Context, session *entities.Session) error
	DeleteByToken(ctx context.Context, refToken string) error
}
