package usecases_test

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/SemgaTeam/semga-stream/internal/core/usecases"
	"github.com/google/uuid"

	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const maxAudioSize = 50 * 1024 * 1024

func newTestConfig() *config.Config {
	return &config.Config{
		Limits: config.Limits{
			MaxAudioSize: maxAudioSize,
		},
		AllowedExtensions: map[string]struct{}{
			".mp3": {},
			".wav": {},
			".ogg": {},
		},
		AllowedMimeTypes: map[string]struct{}{
			"audio/mpeg": {},
			"audio/wav":  {},
			"audio/ogg":  {},
		},
		UploadPath: "uploads",
	}
}

func makeFileHeader(filename string, size int64, mimeType string) *multipart.FileHeader {
	header := textproto.MIMEHeader{}
	header.Set("Content-Type", mimeType)
	return &multipart.FileHeader{
		Filename: filename,
		Size:     size,
		Header:   header,
	}
}

func TestValidateFileUseCase_ValidAudio(t *testing.T) {
	uc := usecases.NewValidateFileUseCase(newTestConfig())
	fileHeader := makeFileHeader("track.WAV", maxAudioSize-1, "audio/wav")

	if err := uc.Execute(context.Background(), fileHeader); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateFileUseCase_InvalidExtension(t *testing.T) {
	uc := usecases.NewValidateFileUseCase(newTestConfig())
	fileHeader := makeFileHeader("track.flac", maxAudioSize-1, "audio/wav")

	if err := uc.Execute(context.Background(), fileHeader); err != e.ErrInvalidExtension {
		t.Fatalf("expected invalid extension error, got %v", err)
	}
}

func TestValidateFileUseCase_TooBig(t *testing.T) {
	uc := usecases.NewValidateFileUseCase(newTestConfig())
	fileHeader := makeFileHeader("track.mp3", maxAudioSize+1, "audio/mpeg")

	if err := uc.Execute(context.Background(), fileHeader); err != e.ErrFileTooBig {
		t.Fatalf("expected file too big error, got %v", err)
	}
}

func TestValidateFileUseCase_InvalidMIME(t *testing.T) {
	uc := usecases.NewValidateFileUseCase(newTestConfig())
	fileHeader := makeFileHeader("track.mp3", maxAudioSize-1, "audio/aac")

	if err := uc.Execute(context.Background(), fileHeader); err != e.ErrInvalidMIMEType {
		t.Fatalf("expected invalid mime type error, got %v", err)
	}
}

func TestSaveFileUseCase_SavesPendingMediaFile(t *testing.T) {
	conf := newTestConfig()
	storage := &fakeStorage{generateURL: "http://upload.local"}
	mediaRepo := &fakeMediaFile{}
	analyzer := &fakeAudioAnalyzer{duration: 3 * time.Second}
	uc := usecases.NewSaveFileUseCase(conf, storage, mediaRepo, analyzer)

	ownerID := uuid.New()
	fileHeader := makeFileHeader("song.MP3", 12345, "audio/mpeg")

	uploadURL, mediaFile, err := uc.Execute(context.Background(), fileHeader, ownerID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if uploadURL != storage.generateURL {
		t.Fatalf("unexpected upload URL: %s", uploadURL)
	}
	if mediaFile == nil {
		t.Fatalf("expected media file, got nil")
	}

	pathPrefix := fmt.Sprintf("%s/%s/", conf.UploadPath, ownerID.String())
	if !strings.HasPrefix(mediaFile.FilePath, pathPrefix) {
		t.Fatalf("unexpected file path: %s", mediaFile.FilePath)
	}
	uniqueName := strings.TrimPrefix(mediaFile.FilePath, pathPrefix)
	if uniqueName == "" {
		t.Fatalf("expected unique name in path")
	}
	ext := strings.ToLower(filepath.Ext(uniqueName))
	if ext != ".mp3" {
		t.Fatalf("expected .mp3 extension, got %s", ext)
	}
	if uniqueName == fileHeader.Filename {
		t.Fatalf("expected unique name to differ from original filename")
	}
	if storage.lastGeneratePath != mediaFile.FilePath {
		t.Fatalf("storage.GenerateUploadURL called with %s, expected %s", storage.lastGeneratePath, mediaFile.FilePath)
	}

	if mediaFile.OwnerID != ownerID {
		t.Fatalf("unexpected owner ID")
	}
	if mediaFile.Type != "audio" {
		t.Fatalf("unexpected type: %s", mediaFile.Type)
	}
	if mediaFile.FileName != fileHeader.Filename {
		t.Fatalf("unexpected file name: %s", mediaFile.FileName)
	}
	if mediaFile.FileSize != fileHeader.Size {
		t.Fatalf("unexpected file size: %d", mediaFile.FileSize)
	}
	if mediaFile.MimeType != fileHeader.Header.Get("Content-Type") {
		t.Fatalf("unexpected mime type: %s", mediaFile.MimeType)
	}
	if mediaFile.DurationSeconds != 3 {
		t.Fatalf("unexpected duration: %d", mediaFile.DurationSeconds)
	}
	if mediaFile.Status != entities.StatusPending {
		t.Fatalf("expected status pending, got %s", mediaFile.Status)
	}
	if mediaFile.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
}

func TestGetUserFilesUseCase_ReturnsOnlyUserFiles(t *testing.T) {
	conf := newTestConfig()
	userID := uuid.New()
	mediaRepo := &fakeMediaFile{
		byUser: map[uuid.UUID][]entities.MediaFile{
			userID: {
				{ID: uuid.New(), OwnerID: userID},
				{ID: uuid.New(), OwnerID: userID},
			},
		},
	}
	uc := usecases.NewGetUserFilesUseCase(conf, mediaRepo)

	files, err := uc.Execute(context.Background(), userID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	for _, f := range files {
		if f.OwnerID != userID {
			t.Fatalf("expected file owner %s, got %s", userID, f.OwnerID)
		}
	}
}

func TestDeleteFileUseCase_DeletesFile(t *testing.T) {
	conf := newTestConfig()
	fileID := uuid.New()
	originalPath := "uploads/u/file.mp3"
	mediaRepo := &fakeMediaFile{
		byID: map[uuid.UUID]*entities.MediaFile{
			fileID: {
				ID:       fileID,
				FilePath: originalPath,
			},
		},
	}
	storage := &fakeStorage{}
	uc := usecases.NewDeleteFileUseCase(conf, mediaRepo, storage)

	if err := uc.Execute(context.Background(), fileID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(mediaRepo.saved) == 0 {
		t.Fatalf("expected media file to be saved")
	}
	saved := mediaRepo.saved[len(mediaRepo.saved)-1]
	if saved.FilePath != "" {
		t.Fatalf("expected file path to be empty, got %s", saved.FilePath)
	}
	if !saved.IsDeleted {
		t.Fatalf("expected file to be marked deleted")
	}
	if storage.lastDeletePath != originalPath {
		t.Fatalf("expected storage.Delete to be called with %s, got %s", originalPath, storage.lastDeletePath)
	}
}

func TestCompleteUploadUseCase_MarksUploaded(t *testing.T) {
	conf := newTestConfig()
	fileID := uuid.New()
	ownerID := uuid.New()
	path := "uploads/u/file.mp3"
	mediaRepo := &fakeMediaFile{
		byID: map[uuid.UUID]*entities.MediaFile{
			fileID: {
				ID:       fileID,
				OwnerID:  ownerID,
				FilePath: path,
				Status:   entities.StatusPending,
			},
		},
	}
	storage := &fakeStorage{
		fileExists: map[string]bool{
			path: true,
		},
	}
	uc := usecases.NewCompleteUploadUseCase(conf, storage, mediaRepo)

	if err := uc.Execute(context.Background(), ownerID, fileID); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mediaRepo.saved) == 0 {
		t.Fatalf("expected media file to be saved")
	}
	saved := mediaRepo.saved[len(mediaRepo.saved)-1]
	if saved.Status != entities.StatusUploaded {
		t.Fatalf("expected status uploaded, got %s", saved.Status)
	}
	if storage.lastExistsPath != path {
		t.Fatalf("expected FileExists to check %s, got %s", path, storage.lastExistsPath)
	}
}
