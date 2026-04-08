package usecases

import (
	"context"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
)

type RegisterAccountDTO struct {
	FullName string
	Username string
	Password string
}

type RegisterAccountUsecase struct {
	config      *config.Config
	accountRepo interfaces.AccountRepository
	hasher      interfaces.IPasswordHasher
}

func NewRegisterAccountUsecase(c *config.Config, r interfaces.AccountRepository, h interfaces.IPasswordHasher) *RegisterAccountUsecase {
	return &RegisterAccountUsecase{
		config:      c,
		accountRepo: r,
		hasher:      h,
	}
}

func (r *RegisterAccountUsecase) Execute(ctx context.Context, a RegisterAccountDTO) (string, error) {
	exists, err := r.accountRepo.ExistsByUsername(ctx, a.Username)
	if err != nil {
		return "", err
	}
	if exists {
		return "", domainErrors.NewError("invalid username: user already exists")
	}

	hashPassword, err := r.hasher.Hash(a.Password)
	if err != nil {
		return "", err
	}
	roles := []entities.Role{entities.RoleUser}

	// Sent id = "" as account is new! Sent createdAt := time.Now() cause we just register a new account
	newAccount, err := entities.NewAccount(uuid.Nil, a.FullName, a.Username, hashPassword, roles, time.Now())
	if err != nil {
		return "", err
	}

	ID, err := r.accountRepo.Save(ctx, newAccount)
	if err != nil {
		return "", err
	}

	return ID, nil
}
