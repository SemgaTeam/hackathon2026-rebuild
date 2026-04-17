//go:build e2e

package test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	"github.com/SemgaTeam/semga-stream/internal/core/usecases"
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/db"
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestUseCasesE2E_FileLifecycle(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conf := loadE2EConfig(t)
	runMigrationsForE2E(t, conf)

	pool, err := db.InitDB(ctx, conf.Postgres.URL)
	if err != nil {
		t.Fatalf("init db: %v", err)
	}

	s3Client := newE2ES3Client(t, conf)
	ensureBucketExists(t, ctx, conf, s3Client)
	cleanupE2EState(t, ctx, pool, conf, s3Client)
	t.Cleanup(pool.Close)
	t.Cleanup(func() {
		cleanupE2EState(t, context.Background(), pool, conf, s3Client)
	})

	userID := seedUser(t, ctx, pool)
	fileHeader, audioBytes := newAudioFileHeader(t, "track.wav", "audio/wav")

	validateUC := usecases.NewValidateFileUseCase(conf)
	if err := validateUC.Execute(ctx, fileHeader); err != nil {
		t.Fatalf("validate file: %v", err)
	}

	storageRepo, err := repository.NewStorageRepository(conf)
	if err != nil {
		t.Fatalf("new storage repository: %v", err)
	}
	mediaRepo := repository.NewMediaFileRepository(conf, pool)
	audioAnalyzer := repository.NewAudioAnalyzer()

	saveUC := usecases.NewSaveFileUseCase(conf, storageRepo, mediaRepo, audioAnalyzer)
	completeUC := usecases.NewCompleteUploadUseCase(conf, storageRepo, mediaRepo)
	getFilesUC := usecases.NewGetUserFilesUseCase(conf, mediaRepo)
	deleteUC := usecases.NewDeleteFileUseCase(conf, mediaRepo, storageRepo)

	uploadURL, mediaFile, err := saveUC.Execute(ctx, fileHeader, userID)
	if err != nil {
		t.Fatalf("save file: %v", err)
	}
	if mediaFile == nil {
		t.Fatal("expected media file to be created")
	}
	if mediaFile.Status != entities.StatusPending {
		t.Fatalf("expected pending status after save, got %s", mediaFile.Status)
	}

	if err := uploadToPresignedURL(ctx, uploadURL, "audio/wav", audioBytes); err != nil {
		t.Fatalf("upload file to storage: %v", err)
	}

	if err := completeUC.Execute(ctx, userID, mediaFile.ID); err != nil {
		t.Fatalf("complete upload: %v", err)
	}

	files, err := getFilesUC.Execute(ctx, userID)
	if err != nil {
		t.Fatalf("get user files: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Status != entities.StatusUploaded {
		t.Fatalf("expected uploaded status, got %s", files[0].Status)
	}
	if files[0].FilePath == "" {
		t.Fatal("expected stored file path to be set")
	}

	if err := deleteUC.Execute(ctx, mediaFile.ID); err != nil {
		t.Fatalf("delete file: %v", err)
	}

	deletedFile, err := mediaRepo.ByID(ctx, mediaFile.ID)
	if err != nil {
		t.Fatalf("load deleted file: %v", err)
	}
	if deletedFile == nil {
		t.Fatal("expected deleted file record to remain in db")
	}
	if !deletedFile.IsDeleted {
		t.Fatal("expected file to be marked deleted")
	}
	if deletedFile.FilePath != "" {
		t.Fatalf("expected file path to be cleared, got %q", deletedFile.FilePath)
	}

	exists, err := storageRepo.FileExists(ctx, mediaFile.FilePath)
	if err != nil {
		t.Fatalf("check file existence after delete: %v", err)
	}
	if exists {
		t.Fatalf("expected storage object %q to be deleted", mediaFile.FilePath)
	}
}

func loadE2EConfig(t *testing.T) *config.Config {
	t.Helper()

	defaults := map[string]string{
		"POSTGRES_URL":                  "postgresql://test_user:test_password@localhost:55432/semga_test?sslmode=disable",
		"SIGNING_KEY":                   "test-signing-key",
		"S3_URL":                        "http://localhost:59000",
		"S3_BUCKET_NAME":                "semga-test",
		"S3_REGION":                     "us-east-1",
		"S3_ACCESS_KEY_ID":              "GK1834094781786f8dde242381",
		"S3_SECRET_ACCESS_KEY":          "6cb5fe16ca3df92f3c6700de488fd90d4b84802a6e89e5da7445a9274d23765d",
		"S3_PRESIGN_EXPIRATION_SECONDS": "300",
		"S3_UPLOAD_PATH":                "test-uploads",
	}

	for key, value := range defaults {
		if os.Getenv(key) == "" {
			t.Setenv(key, value)
		}
	}

	conf, err := config.GetConfig()
	if err != nil {
		t.Fatalf("load e2e config: %v", err)
	}

	return conf
}

func runMigrationsForE2E(t *testing.T, conf *config.Config) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}

	migrationsPath := filepath.Clean(filepath.Join(wd, "..", "migrations"))
	if err := db.RunMigrations(&conf.Postgres, migrationsPath); err != nil {
		t.Fatalf("run migrations: %v", err)
	}
}

func seedUser(t *testing.T, ctx context.Context, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	userID := uuid.New()
	_, err := pool.Exec(ctx,
		`INSERT INTO users (id, login, full_name, password_hash)
		 VALUES ($1, $2, $3, $4)`,
		userID,
		fmt.Sprintf("e2e-%s", strings.ToLower(userID.String()[:8])),
		"E2E Test User",
		"test-password-hash",
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	return userID
}

func cleanupE2EState(t *testing.T, ctx context.Context, pool *pgxpool.Pool, conf *config.Config, client *s3.Client) {
	t.Helper()

	if _, err := pool.Exec(ctx, `TRUNCATE TABLE media_files, users RESTART IDENTITY CASCADE`); err != nil {
		t.Fatalf("truncate test tables: %v", err)
	}

	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(conf.Storage.Bucket),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			t.Fatalf("list test objects: %v", err)
		}

		for _, obj := range page.Contents {
			if obj.Key == nil {
				continue
			}

			if _, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: aws.String(conf.Storage.Bucket),
				Key:    obj.Key,
			}); err != nil {
				t.Fatalf("delete test object %q: %v", aws.ToString(obj.Key), err)
			}
		}
	}
}

func newE2ES3Client(t *testing.T, conf *config.Config) *s3.Client {
	t.Helper()

	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(conf.Storage.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.Storage.AccessKeyID, conf.Storage.SecretAccessKey, "")),
	)
	if err != nil {
		t.Fatalf("load aws config: %v", err)
	}

	return s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(conf.Storage.URL)
	})
}

func ensureBucketExists(t *testing.T, ctx context.Context, conf *config.Config, client *s3.Client) {
	t.Helper()

	_, err := client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(conf.Storage.Bucket),
	})
	if err == nil || strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") || strings.Contains(err.Error(), "BucketAlreadyExists") {
		return
	}

	t.Fatalf("create bucket %q: %v", conf.Storage.Bucket, err)
}

func newAudioFileHeader(t *testing.T, filename, contentType string) (*multipart.FileHeader, []byte) {
	t.Helper()

	audioBytes := silentWAV(1, 8000)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create multipart file: %v", err)
	}
	if _, err := part.Write(audioBytes); err != nil {
		t.Fatalf("write multipart payload: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "http://example.local/upload", &body)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if err := req.ParseMultipartForm(int64(len(audioBytes) + 1024)); err != nil {
		t.Fatalf("parse multipart form: %v", err)
	}

	fileHeader := req.MultipartForm.File["file"][0]
	fileHeader.Header.Set("Content-Type", contentType)

	return fileHeader, audioBytes
}

func uploadToPresignedURL(ctx context.Context, url, contentType string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("unexpected upload status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	return nil
}

func silentWAV(seconds int, sampleRate int) []byte {
	numSamples := seconds * sampleRate
	dataSize := numSamples * 2
	chunkSize := 36 + dataSize

	buf := &bytes.Buffer{}
	buf.WriteString("RIFF")
	writeLE32(buf, uint32(chunkSize))
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	writeLE32(buf, 16)
	writeLE16(buf, 1)
	writeLE16(buf, 1)
	writeLE32(buf, uint32(sampleRate))
	writeLE32(buf, uint32(sampleRate*2))
	writeLE16(buf, 2)
	writeLE16(buf, 16)
	buf.WriteString("data")
	writeLE32(buf, uint32(dataSize))
	buf.Write(make([]byte, dataSize))

	return buf.Bytes()
}

func writeLE16(buf *bytes.Buffer, value uint16) {
	buf.WriteByte(byte(value))
	buf.WriteByte(byte(value >> 8))
}

func writeLE32(buf *bytes.Buffer, value uint32) {
	buf.WriteByte(byte(value))
	buf.WriteByte(byte(value >> 8))
	buf.WriteByte(byte(value >> 16))
	buf.WriteByte(byte(value >> 24))
}
