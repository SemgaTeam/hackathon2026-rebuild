package main

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/usecases"
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/db"
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/http"
	"github.com/SemgaTeam/semga-stream/internal/infrastructure/repository"
	"github.com/labstack/echo/v4"

	"context"
	"log"
)

func main() {
	ctx := context.Background()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := db.InitDB(ctx, conf.Postgres.URL)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	storageRepo, err := repository.NewStorageRepository(conf)
	if err != nil {
		log.Fatal(err)
	}

	mediaFileRepo := repository.NewMediaFileRepository(conf, pool)
	audioAnalyzer := repository.NewAudioAnalyzer()

	validateUC := usecases.NewValidateFileUseCase(conf)
	saveUC := usecases.NewSaveFileUseCase(conf, storageRepo, mediaFileRepo, audioAnalyzer)
	getFilesUC := usecases.NewGetUserFilesUseCase(conf, mediaFileRepo)
	deleteFileUC := usecases.NewDeleteFileUseCase(conf, mediaFileRepo, storageRepo)
	completeUploadUC := usecases.NewCompleteUploadUseCase(conf, storageRepo, mediaFileRepo)

	ctr := http.NewHTTPController(conf, e, validateUC, saveUC, getFilesUC, deleteFileUC, completeUploadUC)
	ctr.SetupHandlers()

	log.Fatal(ctr.Start())
}
