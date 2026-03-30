package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Config struct {
	Limits Limits
	Signing Signing
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
		AllowedExtensions: allowedExtensions,
		AllowedMimeTypes: allowedMimeTypes,
		AllowedOrigins: allowedOrigins,
		AllowedHeaders: allowedHeaders,
		UploadPath: uploadPath,
	}

	return &conf, nil
}
