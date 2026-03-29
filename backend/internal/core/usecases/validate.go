package usecases

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"

	"mime/multipart"
	"path/filepath"
	"strings"
	"context"
)

type ValidateFileUseCase struct {
	conf *config.Config
}

func NewValidateFileUseCase(conf *config.Config) *ValidateFileUseCase {
	return &ValidateFileUseCase{
		conf,
	}
}

func (uc *ValidateFileUseCase) Execute(ctx context.Context, file multipart.File, fileHeader multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if _, ok := uc.conf.AllowedExtensions[ext]; !ok {
		return e.ErrInvalidExtension
	}

	if fileHeader.Size > uc.conf.Limits.MaxAudioSize {
		return e.ErrFileTooBig
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if _, ok := uc.conf.AllowedMimeTypes[mimeType]; !ok {
		return e.ErrInvalidMIMEType
	}

	return nil
}
