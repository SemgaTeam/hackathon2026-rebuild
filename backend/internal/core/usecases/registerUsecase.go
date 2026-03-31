package usecases

import (
	"context"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/SemgaTeam/semga-stream/internal/core/interfaces"

	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
)

type RegisterAccountDTO struct {
	FullName string
	Username string
	Password string
}

type RegisterAccountUsecase struct {
	accountRepo interfaces.AccountRepository
	hasher      interfaces.PasswordHasher
}

func NewRegisterAccountUsecase(r interfaces.AccountRepository, h interfaces.PasswordHasher) *RegisterAccountUsecase {
	return &RegisterAccountUsecase{
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
	newAccount, err := entities.NewAccount("", a.FullName, a.Username, hashPassword, roles, time.Now())
	if err != nil {
		return "", err
	}

	ID, err := r.accountRepo.Save(ctx, newAccount)
	if err != nil {
		return "", err
	}

	return ID, nil
}
