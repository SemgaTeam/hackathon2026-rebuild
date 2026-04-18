package entities

import (
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/google/uuid"

	"slices"
)

type PlaylistItem struct {
	PlaylistID  uuid.UUID
	MediaFileID uuid.UUID
	Position    int
}

type PlaylistItems struct {
	playlistID uuid.UUID
	items []PlaylistItem
}

func NewPlaylistItems(playlistID uuid.UUID) *PlaylistItems {
	return &PlaylistItems{
		playlistID: playlistID,
		items: []PlaylistItem{},
	}
}

func (p *PlaylistItems) AddToEnd(i PlaylistItem) error {
	maxPos := -1

	for range p.items {
		maxPos += 1
	}

	i.Position = maxPos + 1
	p.items = append(p.items, i)
	
	return nil	
}

func (p *PlaylistItems) Move(from, to int) error {
	items := p.items

	if from < 0 || from >= len(items) {
		return e.ErrInvalidMoveRange
	}
	if to < 0 || to >= len(items) {
		return e.ErrInvalidMoveRange
	}

	if from == to {
		return nil
	}

	moved := items[from]

	if from < to {
		copy(items[from:to], items[from+1:to+1])
	} else {
		copy(items[to+1:from+1], items[to:from])
	}

	items[to] = moved

	for i := range items {
		items[i].Position = i
	}

	return nil
}

func (p *PlaylistItems) SortByPosition() {
	slices.SortFunc(p.items, func(a, b PlaylistItem) int {
		return a.Position - b.Position
	})
}
