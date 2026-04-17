package entities

import (
	"github.com/google/uuid"

	"time"
)

type Playlist struct {
	ID uuid.UUID
	OwnerID uuid.UUID
	Name string
	CreatedAt time.Time
	IsDeleted bool
	Tracks []MediaFile
}
