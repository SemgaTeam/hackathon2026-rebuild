package interfaces

import (
	"context"
)

type IStorage interface {
	FileExists(ctx context.Context, path string) (bool, error)
	GenerateUploadURL(ctx context.Context, path string) (string, error)
	Delete(ctx context.Context, path string) error
}
