package interfaces

import (
	"github.com/SemgaTeam/semga-stream/internal/core/entities"

	"mime/multipart"
	"time"
)

type IMediaFile interface {
	Save(media *entities.MediaFile) error
	GetDuration(file multipart.File, fileHeader multipart.FileHeader) (*time.Duration, error)
}
