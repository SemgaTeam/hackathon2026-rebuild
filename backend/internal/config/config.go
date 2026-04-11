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
	dsn, err := MustEnvString("POSTGRES_URL")
	if err != nil {
		return nil, err
	}

	audioSize, err := EnvWithDefaultInt64("MAX_AUDIO_SIZE", 50 * 1024 * 1024)
	if err != nil {
		return nil, err
	}

	videoSize, err := EnvWithDefaultInt64("MAX_VIDEO_SIZE", 100 * 1024 * 1024)
	if err != nil {
		return nil, err
	}

	signingKey, err := MustEnvString("SIGNING_KEY")
	if err != nil {
		return nil, err
	}

	s3Url, err := MustEnvString("S3_URL")
	if err != nil {
		return nil, err
	}

	bucket, err := MustEnvString("S3_BUCKET_NAME")
	if err != nil {
		return nil, err
	}

	region, err := MustEnvString("S3_REGION")
	if err != nil {
		return nil, err
	}

	accessKeyID, err := MustEnvString("S3_ACCESS_KEY_ID")
	if err != nil {
		return nil, err
	}

	secretAccessKey, err := MustEnvString("S3_SECRET_ACCESS_KEY")
	if err != nil {
		return nil, err
	}

	presignExp, err := EnvWithDefaultInt("S3_PRESIGN_EXPIRATION_SECONDS", 5 * 60)
	if err != nil {
		return nil, err
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

	uploadPath := EnvWithDefaultString("S3_UPLOAD_PATH", "uploads")


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

func MustEnvString(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", errors.New(name + " is not set")
	}
	return value, nil
}

func EnvWithDefaultInt64(name string, defaultValue int64) (int64, error) {
	var value = defaultValue
	var err error
	valueStr := os.Getenv(name)
	if valueStr != "" {
		value, err = strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			return 0, errors.New(name + " is invalid")
		}
	}

	return value, nil
}

func EnvWithDefaultInt(name string, defaultValue int) (int, error) {
	var value = defaultValue
	var err error
	valueStr := os.Getenv(name)
	if valueStr != "" {
		value, err = strconv.Atoi(valueStr)
		if err != nil {
			return 0, errors.New(name + " is invalid")
		}
	}

	return value, nil
}

func EnvWithDefaultString(name string, defaultValue string) string {
	var value = defaultValue

	valueStr := os.Getenv(name)
	if valueStr != "" {
		value = valueStr
	}

	return value
}
