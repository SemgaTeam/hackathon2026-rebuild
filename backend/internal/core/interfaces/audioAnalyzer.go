package interfaces

import (
	"context"
	"mime/multipart"
	"time"
)

type IAudioAnalyzer interface {
	GetDuration(ctx context.Context, fileHeader *multipart.FileHeader) (*time.Duration, error)
}
