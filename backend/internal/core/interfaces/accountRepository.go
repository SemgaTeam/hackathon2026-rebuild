package interfaces

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type AccountRepository interface {
	Save(ctx context.Context, account *entities.Account) (string, error)
	FindByID(ctx context.Context, id string) (*entities.Account, error)
	FindByUsername(ctx context.Context, username string) (*entities.Account, error)
	FindAll(ctx context.Context) ([]entities.Account, error)
	UpdateAccount(ctx context.Context, account *entities.Account) error
	DeleteAccount(ctx context.Context, id string) error
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}
