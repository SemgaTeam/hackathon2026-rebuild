package interfaces

import "context"

type IStorage interface {
	GenerateUploadURL(ctx context.Context, path string) (string, error)
}
