package config

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"errors"
	"os"
	"strconv"
)

type Config struct {
	Limits Limits
	Signing Signing
	Storage Storage
	Postgres Postgres
	AllowedExtensions map[string]struct{}
	AllowedMimeTypes map[string]struct{}
	AllowedOrigins []string
	AllowedHeaders []string
	UploadPath string
}

type (
	Limits struct {
		MaxAudioSize int64
		MaxVideoSize int64
	}
	
	Signing struct {
		Key string
		Method jwt.SigningMethod
	}

	Storage struct {
		URL string
		Bucket string
		AccessKeyID string
		SecretAccessKey string
		Region string
		PresignExpirationSeconds int
	}

	Postgres struct {
		URL string
	}
)

func GetConfig() (*Config, error) {
	dsn := os.Getenv("POSTGRES_URL")
	if dsn == "" {
		return nil, errors.New("POSTGRES_URL is not set")
	}

	audioSizeStr := os.Getenv("MAX_AUDIO_SIZE")
	if audioSizeStr == "" {
		return nil, errors.New("MAX_AUDIO_SIZE is not set")
	}
	audioSize, err := strconv.ParseInt(audioSizeStr, 10, 64)
	if err != nil {
		return nil, errors.New("MAX_AUDIO_SIZE is invalid")
	}

	videoSizeStr := os.Getenv("MAX_VIDEO_SIZE")
	if videoSizeStr == "" {
		return nil, errors.New("MAX_VIDEO_SIZE is not set")
	}

	videoSize, err := strconv.ParseInt(videoSizeStr, 10, 64)
	if err != nil {
		return nil, errors.New("MAX_VIDEO_SIZE is invalid")
	}

	signingKey := os.Getenv("SIGNING_KEY")
	if signingKey == "" {
		return nil, errors.New("empty SIGNING_KEY")
	}

	s3Url := os.Getenv("S3_URL")
	if s3Url == "" {
		return nil, errors.New("S3_URL is not set")
	}

	bucket := os.Getenv("S3_BUCKET_NAME")
	if bucket == "" {
		return nil, errors.New("S3_BUCKET_NAME is not set")
	}

	region := os.Getenv("S3_REGION")
	if region == "" {
		return nil, errors.New("S3_REGION is not set")
	}

	accessKeyID := os.Getenv("S3_ACCESS_KEY_ID")
	if accessKeyID == "" {
		return nil, errors.New("S3_ACCESS_KEY_ID is not set")
	}

	secretAccessKey := os.Getenv("S3_SECRET_ACCESS_KEY")
	if secretAccessKey == "" {
		return nil, errors.New("S3_SECRET_ACCESS_KEY is not set")
	}

	presignExpStr := os.Getenv("S3_PRESIGN_EXPIRATION_SECONDS")
	if presignExpStr == "" {
		return nil, errors.New("S3_PRESIGN_EXPIRATION_SECONDS is not set")
	}

	presignExp, err := strconv.Atoi(presignExpStr)
	if err != nil {
		return nil, errors.New("S3_PRESIGN_EXPIRATION_SECONDS is invalid")
	}

	signingMethod := jwt.SigningMethodHS256

	allowedExtensions := map[string]struct{} {
		".mp3": {},
		".wav": {},
		".ogg": {},
	}
	allowedMimeTypes := map[string]struct{} {
		"audio/mpeg": {},
		"audio/wav": {},
		"audio/ogg": {},
	}

	allowedOrigins := []string {"http://localhost:5173"} // TODO: get in env variable
	allowedHeaders := []string {echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}

	uploadPath := "uploads"


	conf := Config{
		Limits: Limits{
			MaxAudioSize: audioSize,
			MaxVideoSize: videoSize,
		},
		Signing: Signing{
			Key: signingKey,
			Method: signingMethod,
		},
		Storage: Storage{
			URL: s3Url,
			Bucket: bucket,
			AccessKeyID: accessKeyID,
			SecretAccessKey: secretAccessKey,
			Region: region,
			PresignExpirationSeconds: presignExp,
		},
		Postgres: Postgres{
			URL: dsn,
		},
		AllowedExtensions: allowedExtensions,
		AllowedMimeTypes: allowedMimeTypes,
		AllowedOrigins: allowedOrigins,
		AllowedHeaders: allowedHeaders,
		UploadPath: uploadPath,
	}

	return &conf, nil
}
