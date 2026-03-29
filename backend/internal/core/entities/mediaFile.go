package entities

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	Type            string
	FileName        string
	FilePath        string
	FileSize        int64
	MimeType        string
	DurationSeconds int
	CreatedAt       time.Time
	IsDeleted       bool
}
