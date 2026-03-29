package config

import (
	"errors"
	"os"
	"strconv"
)

type Config struct {
	Limits Limits
	AllowedExtensions AllowedExtensions
}

type (
	Limits struct {
		MaxAudioSize int64
		MaxVideoSize int64
	}

	AllowedExtensions map[string]struct{}
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

	var allowedExtensions AllowedExtensions = map[string]struct{} {
		".mp3": {},
		".wav": {},
		".ogg": {},
	}

	conf := Config{
		Limits: Limits{
			MaxAudioSize: audioSize,
			MaxVideoSize: videoSize,
		},
		AllowedExtensions: allowedExtensions,
	}

	return &conf, nil
}
