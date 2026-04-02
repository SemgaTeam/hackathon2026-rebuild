package interfaces

import (
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/google/uuid"

	"context"
)

type IMediaFile interface {
	Save(ctx context.Context, media *entities.MediaFile) error
	ByID(ctx context.Context, id uuid.UUID) (*entities.MediaFile, error)
	ByUserID(ctx context.Context, userId uuid.UUID) ([]entities.MediaFile, error)
	ByPath(ctx context.Context, path string) (*entities.MediaFile, error)
}
