package interfaces

import (
	"github.com/SemgaTeam/semga-stream/internal/core/dto"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type IToken interface {
	GenerateTokens(account *entities.Account) (dto.Tokens, error)
	VerifyRefresh(refToken string) (string, error)
	GenerateAccess(account *entities.Account) (string error)
}
