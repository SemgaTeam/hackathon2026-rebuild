package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"context"
)

type CompleteUploadUseCase struct {
	conf      *config.Config
	storage   i.IStorage
	mediaFile i.IMediaFile
}

func NewCompleteUploadUseCase(conf *config.Config, storage i.IStorage, mediaFile i.IMediaFile) *CompleteUploadUseCase {
	return &CompleteUploadUseCase{
		conf,
		storage,
		mediaFile,
	}
}

func (uc *CompleteUploadUseCase) Execute(ctx context.Context, ownerId uuid.UUID, fileId uuid.UUID) error {
	file, err := uc.mediaFile.ByID(ctx, fileId)
	if err != nil {
		return err
	}
	if file == nil {
		return e.ErrFileNotFound
	}

	exists, err := uc.storage.FileExists(ctx, file.FilePath)
	if err != nil {
		return err
	}

	if !exists {
		return e.ErrFileNotFound
	}

	file.Status = entities.StatusUploaded
	if err := uc.mediaFile.Save(ctx, file); err != nil {
		return err
	}

	return nil
}
