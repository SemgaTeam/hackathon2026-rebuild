package interfaces

import "github.com/SemgaTeam/semga-stream/internal/core/entities"

type AccountRepository interface {
	Save(account *entities.Account) (string, error)
	FindByID(id string) (*entities.Account, error)
	FindByUsername(username string) (*entities.Account, error)
	FindAll() ([]entities.Account, error)
	UpdateAccount(account *entities.Account) error
	DeleteAccount(id string) error
	ExistsByUsername(username string) bool
}
