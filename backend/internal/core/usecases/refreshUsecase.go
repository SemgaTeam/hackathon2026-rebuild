package usecases

import (
	"context"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/config"
	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"
)

type RefreshTokenUsecase struct {
	config       *config.Config
	tokenService interfaces.IToken
	accountRepo  interfaces.AccountRepository
	sessionRepo  interfaces.SessionRepository
}

func NewRefreshTokenUsecase(c *config.Config, t interfaces.IToken, a interfaces.AccountRepository, s interfaces.SessionRepository) *RefreshTokenUsecase {
	return &RefreshTokenUsecase{
		config:       c,
		tokenService: t,
		accountRepo:  a,
		sessionRepo:  s,
	}
}

func (r *RefreshTokenUsecase) Execute(ctx context.Context, refToken string) (string, error) {
	if refToken == "" {
		return "", domainErrors.NewError("invalid refresh token: token must be not empty")
	}

	accountID, err := r.tokenService.VerifyRefresh(ctx, refToken)
	if err != nil {
		return "", domainErrors.NewError("invalid refresh token: verification failed")
	}

	session, err := r.sessionRepo.FindByToken(ctx, refToken)
	if err != nil {
		return "", domainErrors.NewError("invalid session: there is no session")
	}

	if time.Now().After(session.ExpiresAt) {
		return "", domainErrors.NewError("invalid session: session expired")
	}

	account, err := r.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return "", domainErrors.NewError("invalid session: account was removed")
	}

	accessToken, err := r.tokenService.GenerateAccess(ctx, account)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
