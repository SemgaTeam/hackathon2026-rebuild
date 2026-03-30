package http

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	uc "github.com/SemgaTeam/semga-stream/internal/core/usecases"
	e "github.com/SemgaTeam/semga-stream/internal/infrastructure/http/errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echojwt "github.com/labstack/echo-jwt/v4"

	"net/http"
	"path/filepath"
	"strings"
)

type Controller struct {
	conf *config.Config
	e *echo.Echo
	validateUC *uc.ValidateFileUseCase
	saveUC *uc.SaveFileUseCase
}

func NewHTTPController(conf *config.Config, e *echo.Echo, validateUC *uc.ValidateFileUseCase, saveUC *uc.SaveFileUseCase) *Controller {
	return &Controller{
		conf,
		e,
		validateUC,
		saveUC,
	}
}

func (ctr *Controller) SetupHandlers() {
	api := ctr.e.Group("/api")

	ctr.e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(ctr.conf.Signing.Key),
		TokenLookup:   "cookie:access_token",
		ContextKey:    "access_token",
		SigningMethod: ctr.conf.Signing.Method.Alg(),
	}))
	ctr.e.Use(middleware.AddTrailingSlash())
	ctr.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: ctr.conf.AllowedOrigins,
		AllowHeaders: ctr.conf.AllowedHeaders,
	}))

	ctr.e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Request().Body = http.MaxBytesReader(
				c.Response(),
				c.Request().Body,
				ctr.conf.Limits.MaxAudioSize,
			)
			return next(c)
		}
	}) 

	api.POST("/upload", ctr.UploadHandler)
}

func (ctr *Controller) UploadHandler(c echo.Context) error {
	fileHeader, err := c.FormFile("file")	
	if err != nil {
		return e.BadRequest("file not provided")
	}

	if fileHeader.Size > ctr.conf.Limits.MaxAudioSize {
		return e.BadRequest("file too large")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if _, ok := ctr.conf.AllowedExtensions[ext]; !ok {
		return e.BadRequest("invalid file extension")
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if _, ok := ctr.conf.AllowedMimeTypes[mimeType]; !ok {
		return e.BadRequest("invalid MIME type")
	}

	ctx := c.Request().Context()

	var ownerId uuid.UUID
	uploadUrl, _, err := ctr.saveUC.Execute(ctx, fileHeader, ownerId) 
	if err != nil {
		return err
	}

	response := map[string]string{
		"upload_url": uploadUrl,
	}

	return c.JSON(http.StatusOK, response)
}
