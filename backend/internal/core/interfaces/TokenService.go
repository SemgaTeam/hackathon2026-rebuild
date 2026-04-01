package interfaces

import "github.com/SemgaTeam/semga-stream/internal/core/entities"

type TokenNRefToken struct {
	AccessToken  string
	RefreshToken string
}

type TokenService interface {
	GenerateTokenNRefToken(account *entities.Account) (TokenNRefToken, error)
	VerifyRefToken(refToken string) (string, error)
}
