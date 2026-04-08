package entities

import (
	"regexp"
	"slices"
	"time"

	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/google/uuid"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	fullNameRegex = regexp.MustCompile(`^[А-Яа-яЁё]+\s+[А-Яа-яЁё]+\s+[А-Яа-яЁё]+$`)
)

type Account struct {
	ID           uuid.UUID
	FullName     string
	Username     string
	PasswordHash string
	Roles        []Role

	CreatedAt time.Time
}

// Empty id if just created, not from ps, id is filled if it is from persistency storage
// Creating time defined by client
func NewAccount(id uuid.UUID, fullName, username, passwordHash string, roles []Role, createdAt time.Time) (*Account, error) {

	if !usernameRegex.MatchString(username) {
		return nil, domainErrors.NewError("invalid username: only latin letters and digits are allowed")
	}

	if !fullNameRegex.MatchString(fullName) {
		return nil, domainErrors.NewError("invalid fullName: only russian letters are allowed, only three words separated by spaces")
	}

	if len(roles) == 0 {
		return nil, domainErrors.NewError("invalid roles: at least one role must be specified")
	}

	if passwordHash == "" {
		return nil, domainErrors.NewError("invalid passwordHash: cannot be empty")
	}

	if createdAt.Unix() > time.Now().Unix() {
		return nil, domainErrors.NewError("invalid createdAt: cannot be in the future")
	}

	return &Account{
		ID:           id,
		FullName:     fullName,
		Username:     username,
		Roles:        roles,
		PasswordHash: passwordHash,

		CreatedAt: createdAt}, nil
}
func (a *Account) HasRole(role Role) bool {
	return slices.Contains(a.Roles, role)
}

func (a *Account) IsAdmin() bool {
	return a.HasRole("admin")
}
