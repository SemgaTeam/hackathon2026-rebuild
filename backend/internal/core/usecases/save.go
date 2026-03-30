package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	i "github.com/SemgaTeam/semga-stream/internal/core/interfaces"
	"github.com/google/uuid"

	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
	"context"
)

type SaveFileUseCase struct {
	conf *config.Config
	storage i.IStorage
	mediaFile i.IMediaFile
}

func NewSaveFileUseCase(conf *config.Config, storage i.IStorage, mediaFile i.IMediaFile) *SaveFileUseCase {
	return &SaveFileUseCase{
		conf,
		storage,
		mediaFile,
	}
}

func (uc *SaveFileUseCase) Execute(ctx context.Context, file *multipart.File, fileHeader *multipart.FileHeader, ownerId uuid.UUID) (string, *entities.MediaFile, error) {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	uniqueName := uuid.New().String() + ext

	path := uc.conf.UploadPath + uniqueName

	uploadUrl, err := uc.storage.GenerateUploadURL(ctx, path)
	if err != nil {
		return "", nil, err
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	filename := filepath.Base(fileHeader.Filename)
	duration, err := uc.mediaFile.GetDuration(ctx, file, fileHeader)
	if err != nil {
		return "", nil, err
	}

	mediaFile := entities.MediaFile{
		OwnerID: ownerId,
		Type: "audio", // only audio by now
		FileName: filename,
		FilePath: path,
		FileSize: fileHeader.Size,
		MimeType: mimeType,
		DurationSeconds: int(duration.Seconds()),
		CreatedAt: time.Now(),
	}

	if err := uc.mediaFile.Save(ctx, &mediaFile); err != nil {
		return "", nil, err
	}

	return uploadUrl, &mediaFile, nil
}
