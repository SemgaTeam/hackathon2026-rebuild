package entities

import "github.com/google/uuid"

type PlaylistItem struct {
	PlaylistID uuid.UUID
	MediaFileID uuid.UUID
}
