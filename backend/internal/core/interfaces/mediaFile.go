package interfaces

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/google/uuid"

	"mime/multipart"
	"time"
)

type IMediaFile interface {
	Save(ctx context.Context, media *entities.MediaFile) error
	GetDuration(ctx context.Context, fileHeader *multipart.FileHeader) (*time.Duration, error)
	ByID(ctx context.Context, id uuid.UUID) (*entities.MediaFile, error)
	ByUserID(ctx context.Context, userId uuid.UUID) ([]entities.MediaFile, error)
	ByPath(ctx context.Context, path string) (*entities.MediaFile, error)
}
