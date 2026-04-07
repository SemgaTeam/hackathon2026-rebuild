package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"context"
	"fmt"
)

type CompleteUploadUseCase struct {
	conf *config.Config
	storage i.IStorage
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
	path := fmt.Sprintf("%s/%s/%s", uc.conf.UploadPath, ownerId, fileId)
	exists, err := uc.storage.FileExists(ctx, path)	
	if err != nil {
		return err
	}

	if !exists {
		return e.ErrFileNotFound
	}

	file, err := uc.mediaFile.ByPath(ctx, path)
	if err != nil {
		return err
	}

	if file == nil {
		return e.ErrFileNotFound
	}

	file.Status = entities.StatusUploaded
	if err := uc.mediaFile.Save(ctx, file); err != nil {
		return err
	}

	return nil
}
