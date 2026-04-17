package repository

import (
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/vorbis"
	"github.com/faiface/beep/wav"

	"context"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type AudioAnalyzer struct{}

func NewAudioAnalyzer() *AudioAnalyzer {
	return &AudioAnalyzer{}
}

func (r *AudioAnalyzer) GetDuration(ctx context.Context, fileHeader *multipart.FileHeader) (duration *time.Duration, err error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, e.ErrOpeningFile
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			err = e.Unknown(cerr)
		}
	}()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	var streamer beep.StreamSeekCloser
	var format beep.Format

	switch ext {
	case ".mp3":
		streamer, format, err = mp3.Decode(file)
	case ".wav":
		streamer, format, err = wav.Decode(file)
	case ".ogg":
		streamer, format, err = vorbis.Decode(file)
	default:
		return nil, e.ErrInvalidExtension
	}

	if err != nil {
		return nil, e.Unknown(err)
	}

	value := time.Duration(streamer.Len()) * time.Second / time.Duration(format.SampleRate)
	duration = &value

	return
}
