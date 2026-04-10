package usecases_test

import (
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/google/uuid"

	"context"
	"mime/multipart"
	"sync"
	"time"
)

type fakeStorage struct {
	mu sync.Mutex

	fileExists  map[string]bool
	generateURL string

	lastGeneratePath string
	lastDeletePath   string
	lastExistsPath   string
}

func (f *fakeStorage) FileExists(ctx context.Context, path string) (bool, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastExistsPath = path
	return f.fileExists[path], nil
}

func (f *fakeStorage) GenerateUploadURL(ctx context.Context, path string) (string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastGeneratePath = path
	return f.generateURL, nil
}

func (f *fakeStorage) Delete(ctx context.Context, path string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastDeletePath = path
	return nil
}

type fakeMediaFile struct {
	mu sync.Mutex

	byID       map[uuid.UUID]*entities.MediaFile
	byUser     map[uuid.UUID][]entities.MediaFile
	saved      []*entities.MediaFile
	lastByID   uuid.UUID
	lastByUser uuid.UUID
}

func (f *fakeMediaFile) Save(ctx context.Context, media *entities.MediaFile) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	copy := *media
	f.saved = append(f.saved, &copy)
	if f.byID != nil && copy.ID != uuid.Nil {
		f.byID[copy.ID] = &copy
	}
	return nil
}

func (f *fakeMediaFile) ByID(ctx context.Context, id uuid.UUID) (*entities.MediaFile, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastByID = id
	if f.byID == nil {
		return nil, nil
	}
	file := f.byID[id]
	if file == nil {
		return nil, nil
	}
	copy := *file
	return &copy, nil
}

func (f *fakeMediaFile) ByUserID(ctx context.Context, userId uuid.UUID) ([]entities.MediaFile, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.lastByUser = userId
	if f.byUser == nil {
		return nil, nil
	}
	files := f.byUser[userId]
	out := make([]entities.MediaFile, len(files))
	copy(out, files)
	return out, nil
}

func (f *fakeMediaFile) ByPath(ctx context.Context, path string) (*entities.MediaFile, error) {
	return nil, nil
}

type fakeAudioAnalyzer struct {
	duration time.Duration
}

func (f *fakeAudioAnalyzer) GetDuration(ctx context.Context, fileHeader *multipart.FileHeader) (*time.Duration, error) {
	return &f.duration, nil
}
