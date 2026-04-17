package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"context"
)

type PlaylistsUseCase struct {
	conf     *config.Config
	playlist i.IPlaylist	
}

func NewPlaylistsUseCase(conf *config.Config, playlistInterface i.IPlaylist) *PlaylistsUseCase {
	return &PlaylistsUseCase{
		conf: conf,
		playlist: playlistInterface,
	}	
}

func (uc *PlaylistsUseCase) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]entities.Playlist, error) {
	playlists, err := uc.playlist.AllByOwnerID(ctx, userID)	
	if err != nil {
		return nil, err
	}

	for i := range playlists {
		tracks, err := uc.playlist.GetPlaylistTracks(ctx, playlists[i].ID)	
		if err != nil {
			return nil, err
		}

		playlists[i].Tracks = tracks
	}

	return playlists, nil
}

func (uc *PlaylistsUseCase) CreatePlaylist(ctx context.Context, ownerID uuid.UUID, name string) error {
	playlist, err := entities.NewPlaylist(ownerID, name)
	if err != nil {
		return err
	}

	if err := uc.playlist.Save(ctx, playlist); err != nil {
		return err
	}

	return nil
}
