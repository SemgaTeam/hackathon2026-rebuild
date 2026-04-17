package entities

import (
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
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

func NewPlaylist(ownerID uuid.UUID, name string) (*Playlist, error) {
	if ownerID == uuid.Nil {
		return nil, e.ErrInvalidUUID
	}

	if name == "" {
		return nil, e.ErrInvalidName
	}

	playlist := Playlist{
		OwnerID: ownerID,
		Name: name,
	}

	return &playlist, nil
}
