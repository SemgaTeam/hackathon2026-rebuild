package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Limits Limits
	AllowedExtensions map[string]struct{}
	UploadPath string
}

type (
	Limits struct {
		MaxAudioSize int64
		MaxVideoSize int64
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

	allowedExtensions := map[string]struct{} {
		".mp3": {},
		".wav": {},
		".ogg": {},
	}

	uploadPath := "uploads"

	conf := Config{
		Limits: Limits{
			MaxAudioSize: audioSize,
			MaxVideoSize: videoSize,
		},
		AllowedExtensions: allowedExtensions,
		UploadPath: uploadPath,
	}

	return &conf, nil
}
