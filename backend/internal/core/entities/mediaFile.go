package entities

import (
	"time"

	"github.com/google/uuid"
)

type MediaFile struct {
	ID              uuid.UUID `json:"id"`
	OwnerID         uuid.UUID `json:"owner_id"`
	Type            string `json:"type"`
	FileName        string `json:"file_name"`
	FilePath        string `json:"file_path"`
	FileSize        int64 `json:"file_size"`
	MimeType        string `json:"mime_type"`
	DurationSeconds int `json:"duration_seconds"`
	CreatedAt       time.Time `json:"created_at"`
	Status          mediaFileStatus `json:"status"`
	IsDeleted       bool `json:"is_deleted"`
}

type mediaFileStatus string

const StatusPending mediaFileStatus = "pending"
const StatusUploaded mediaFileStatus = "uploaded"
const StatusFailed mediaFileStatus = "failed"
