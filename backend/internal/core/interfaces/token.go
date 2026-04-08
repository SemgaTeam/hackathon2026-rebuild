package interfaces

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/core/dto"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type IToken interface {
	GenerateTokens(ctx context.Context, account *entities.Account) (dto.Tokens, error)
	// Returns account id, error
	VerifyRefresh(ctx context.Context, refToken string) (string, error)
	GenerateAccess(ctx context.Context, account *entities.Account) (string, error)
}
