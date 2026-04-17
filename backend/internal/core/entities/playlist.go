package entities

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	ID uuid.UUID
	OwnerID uuid.UUID
	Name string
	CreatedAt time.Time
	IsDeleted bool
}
