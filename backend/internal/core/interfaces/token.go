package interfaces

import (
	"github.com/SemgaTeam/semga-stream/internal/core/dto"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type IToken interface {
	GenerateTokenNRefToken(account *entities.Account) (dto.Tokens, error)
	VerifyRefToken(refToken string) (string, error)
}
