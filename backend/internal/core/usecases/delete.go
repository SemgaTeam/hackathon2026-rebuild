package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"context"
)

type DeleteFileUseCase struct {
	conf      *config.Config
	mediaFile i.IMediaFile
	storage   i.IStorage
}

func NewDeleteFileUseCase(conf *config.Config, mediaFile i.IMediaFile, storage i.IStorage) *DeleteFileUseCase {
	return &DeleteFileUseCase{
		conf,
		mediaFile,
		storage,
	}
}

func (uc *DeleteFileUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	file, err := uc.mediaFile.ByID(ctx, id)
	if err != nil {
		return err
	}

	path := string([]byte(file.FilePath)) // full copy
	file.FilePath = ""
	file.IsDeleted = true

	if err := uc.mediaFile.Save(ctx, file); err != nil {
		return err
	}

	if err := uc.storage.Delete(ctx, path); err != nil {
		return err
	}

	return nil
}
