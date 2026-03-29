package usecases

import (
	"context"

	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"
)

type GetUserFilesUseCase struct {
	conf *config.Config
	mediaFile i.IMediaFile
}

func NewGetUserFilesUseCase(conf *config.Config, mediaFile i.IMediaFile) *GetUserFilesUseCase {
	return &GetUserFilesUseCase{
		conf,
		mediaFile,
	}
}

func (uc *GetUserFilesUseCase) Execute(ctx context.Context, userId uuid.UUID) ([]entities.MediaFile, error) {
	files, err := uc.mediaFile.ByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	return files, nil
}
