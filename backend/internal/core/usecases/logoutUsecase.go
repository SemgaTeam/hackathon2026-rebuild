package usecases

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/config"
	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"
)

type LogoutUsecase struct {
	config       *config.Config
	sessionRepo  interfaces.SessionRepository
	tokenService interfaces.IToken
}

func NewLogoutUsecase(c *config.Config, s interfaces.SessionRepository, t interfaces.IToken) *LogoutUsecase {
	return &LogoutUsecase{
		config:       c,
		sessionRepo:  s,
		tokenService: t,
	}
}

func (l *LogoutUsecase) Execute(ctx context.Context, refToken string) (string, error) {
	if refToken == "" {
		return "", domainErrors.NewError("invalid refresh token: token must be not empty")
	}

	accountId, err := l.tokenService.VerifyRefresh(ctx, refToken)
	if err != nil {
		return "", domainErrors.NewError("invalid refresh token: verification failed")
	}

	err = l.sessionRepo.DeleteByToken(ctx, refToken)
	if err != nil {
		return "", domainErrors.NewError("invalid session: there is no session")
	}

	return accountId, nil
}
