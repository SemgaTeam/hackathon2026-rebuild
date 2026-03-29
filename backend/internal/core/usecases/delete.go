package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"

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

func (uc *DeleteFileUseCase) Execute(ctx context.Context, path string) error {
	file, err := uc.mediaFile.ByPath(ctx, path)
	if err != nil {
		return err
	}

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
