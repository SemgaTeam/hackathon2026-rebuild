package interfaces

import (
	"time"

	"github.com/SemgaTeam/semga-stream/internal/core/entities"
)

type IMediaFile interface {
	Save(media *entities.MediaFile) error
	GetDuration(path string) (*time.Duration, error)
}
