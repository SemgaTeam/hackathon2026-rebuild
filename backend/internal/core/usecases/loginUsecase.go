package usecases

import (
	"context"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/dto"
	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"
)

type LoginDTO struct {
	Username string
	Password string
}

type LoginUsecase struct {
	config       *config.Config
	accountRepo  interfaces.AccountRepository
	hasher       interfaces.IPasswordHasher
	tokenService interfaces.IToken
	sessionRepo  interfaces.SessionRepository
}

func NewLoginUsecase(c *config.Config, a interfaces.AccountRepository, h interfaces.IPasswordHasher, t interfaces.IToken, s interfaces.SessionRepository) *LoginUsecase {
	return &LoginUsecase{
		config:       c,
		accountRepo:  a,
		hasher:       h,
		tokenService: t,
		sessionRepo:  s,
	}
}

func (l *LoginUsecase) Execute(ctx context.Context, ld LoginDTO) (dto.Tokens, error) {
	account, err := l.accountRepo.FindByUsername(ctx, ld.Username)
	if err != nil {
		return dto.Tokens{}, domainErrors.NewError("credentials error")
	}

	match := l.hasher.Compare(ld.Password, account.PasswordHash)
	if !match {
		return dto.Tokens{}, domainErrors.NewError("credentials error")
	}

	tokens, err := l.tokenService.GenerateTokens(account)
	if err != nil {
		return dto.Tokens{}, err
	}

	exp := time.Now().Add(l.config.RefreshTokenTTL)
	err = l.sessionRepo.Save(ctx, account.ID, tokens.RefreshToken, exp)
	if err != nil {
		return dto.Tokens{}, err
	}

	return tokens, nil
}
