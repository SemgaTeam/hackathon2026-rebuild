package entities

import (
	"regexp"
	"time"

	domainErrors "github.com/SemgaTeam/semga-stream/internal/core/errors"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	fullNameRegex = regexp.MustCompile(`^[А-Яа-яЁё]+\s+[А-Яа-яЁё]+\s+[А-Яа-яЁё]+$`)
)

type Account struct {
	ID           string
	FullName     string
	Username     string
	PasswordHash string
	Roles        []string

	CreatedAt time.Time
}

func NewAccount(id, fullName, username, passwordHash string, roles []string, createdAt time.Time) (*Account, error) {
	if id == "" {
		return nil, domainErrors.NewError("invalid id: cannot be empty")
	}

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
