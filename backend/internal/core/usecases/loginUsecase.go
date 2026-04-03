package usecases

import (
	"context"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/core/dto"
	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"
)

var (
	Duration = 24 * 30 * time.Hour
)

type LoginDTO struct {
	Username string
	Password string
}

type LoginUsecase struct {
	accountRepo  interfaces.AccountRepository
	hasher       interfaces.IPasswordHasher
	tokenService interfaces.IToken
	sessionRepo  interfaces.SessionRepository
}

func NewLoginUsecase(a interfaces.AccountRepository, h interfaces.IPasswordHasher, t interfaces.IToken, s interfaces.SessionRepository) *LoginUsecase {
	return &LoginUsecase{
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

	tokens, err := l.tokenService.GenerateTokenNRefToken(account)
	if err != nil {
		return dto.Tokens{}, err
	}

	exp := time.Now().Add(Duration)
	err = l.sessionRepo.Save(ctx, account.ID, tokens.RefreshToken, exp)
	if err != nil {
		return dto.Tokens{}, err
	}

	return tokens, nil
}
