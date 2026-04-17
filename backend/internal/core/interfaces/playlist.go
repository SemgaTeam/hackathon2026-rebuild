package interfaces

import (
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/google/uuid"

	"context"
)

type IPlaylist interface {
	Save(ctx context.Context, playlist *entities.Playlist) error
	AllByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]entities.Playlist, error)
	ByID(ctx context.Context, id uuid.UUID) (*entities.Playlist, error)
	GetPlaylistTracks(ctx context.Context, playlistID uuid.UUID) ([]entities.MediaFile, error)
}
