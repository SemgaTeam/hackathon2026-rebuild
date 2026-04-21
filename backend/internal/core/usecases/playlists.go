package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"context"
)

type PlaylistsUseCase struct {
	conf      *config.Config
	playlist  i.IPlaylist	
	mediaFile i.IMediaFile
}

func NewPlaylistsUseCase(conf *config.Config, playlistInterface i.IPlaylist, mediaFileInterface i.IMediaFile) *PlaylistsUseCase {
	return &PlaylistsUseCase{
		conf: conf,
		playlist: playlistInterface,
		mediaFile: mediaFileInterface,
	}	
}

func (uc *PlaylistsUseCase) GetUserPlaylists(ctx context.Context, userID uuid.UUID) ([]entities.Playlist, error) {
	playlists, err := uc.playlist.AllByOwnerID(ctx, userID)	
	if err != nil {
		return nil, err
	}

	for i := range playlists {
		if playlists[i].IsDeleted {
			continue
		}

		tracks, err := uc.playlist.GetPlaylistTracks(ctx, playlists[i].ID)	
		if err != nil {
			return nil, err
		}

		playlists[i].Tracks = *tracks
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

func (uc *PlaylistsUseCase) DeletePlaylist(ctx context.Context, playlistID uuid.UUID) error {
	playlist, err := uc.playlist.ByID(ctx, playlistID)
	if err != nil {
		return err
	}

	if playlist == nil {
		return nil
	}

	playlist.IsDeleted = true	

	if err := uc.playlist.Save(ctx, playlist); err != nil {
		return err
	}

	return nil
}

func (uc *PlaylistsUseCase) RenamePlaylist(ctx context.Context, playlistID uuid.UUID, name string) error {
	playlist, err := uc.playlist.ByID(ctx, playlistID)	
	if err != nil {
		return err
	}

	if playlist == nil {
		return e.ErrPlaylistNotFound
	}

	if err = playlist.Update(name); err != nil {
		return err
	}

	if err = uc.playlist.Save(ctx, playlist); err != nil {
		return err
	}

	return nil
}

func (uc *PlaylistsUseCase) AddTrackToEnd(ctx context.Context, playlistID, fileID uuid.UUID) (*entities.Playlist, error) {
	playlist, err := uc.playlist.ByID(ctx, playlistID)	
	if err != nil {
		return nil, err
	}

	if playlist == nil {
		return nil, e.ErrPlaylistNotFound
	}

	tracks, err := uc.playlist.GetPlaylistTracks(ctx, playlist.ID)
	if err != nil {
		return nil, err
	}

	file, err := uc.mediaFile.ByID(ctx, fileID)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, e.ErrFileNotFound
	}

	if playlist.OwnerID != file.OwnerID {
		return nil, e.ErrPlaylistAndFileOwnersNotMatch
	}

	track := entities.PlaylistItem{
		PlaylistID: playlist.ID,
		MediaFileID: file.ID,
	}

	if err := tracks.AddToEnd(track); err != nil {
		return nil, err
	}

	if err := uc.playlist.Save(ctx, playlist); err != nil {
		return nil, err
	}

	return playlist, nil
}

func (uc *PlaylistsUseCase) DeletePlaylistItemByIndex(ctx context.Context, playlistID uuid.UUID, idx int) error {
	playlist, err := uc.playlist.ByID(ctx, playlistID)
	if err != nil {
		return err
	}

	if playlist == nil {
		return e.ErrPlaylistNotFound
	}

	if err = playlist.Tracks.DeleteByIndex(idx); err != nil {
		return err
	}

	if err = uc.playlist.Save(ctx, playlist); err != nil {
		return err
	}

	return nil
}

func (uc *PlaylistsUseCase) MoveItem(ctx context.Context, playlistID uuid.UUID, from, to int) error {
	playlist, err := uc.playlist.ByID(ctx, playlistID)
	if err != nil {
		return err
	}

	if playlist == nil {
		return e.ErrPlaylistNotFound
	}

	if err = playlist.Tracks.Move(from, to); err != nil {
		return err
	}

	if err = uc.playlist.Save(ctx, playlist); err != nil {
		return err
	}

	return nil
}
