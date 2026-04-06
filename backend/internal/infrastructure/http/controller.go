package http

import (
	"errors"

	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	uc "github.com/SemgaTeam/semga-stream/internal/core/usecases"
	e "github.com/SemgaTeam/semga-stream/internal/infrastructure/http/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"net/http"
	"path/filepath"
	"strings"
)

type Controller struct {
	conf         *config.Config
	e            *echo.Echo
	saveUC       *uc.SaveFileUseCase
	getFilesUC   *uc.GetUserFilesUseCase
	deleteFileUC *uc.DeleteFileUseCase
}

func NewHTTPController(conf *config.Config, e *echo.Echo, saveUC *uc.SaveFileUseCase, getFilesUC *uc.GetUserFilesUseCase, deleteFileUC *uc.DeleteFileUseCase) *Controller {
	return &Controller{
		conf,
		e,
		saveUC,
		getFilesUC,
		deleteFileUC,
	}
}

func (ctr *Controller) Start() error {
	return ctr.e.Start(":8080")
}

func (ctr *Controller) SetupHandlers() {
	ctr.e.HTTPErrorHandler = errorHandler
	api := ctr.e.Group("/api")
	me := api.Group("/me")
	files := me.Group("/files")

	ctr.e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(ctr.conf.Signing.Key),
		TokenLookup:   "cookie:access_token",
		ContextKey:    "access_token",
		SigningMethod: ctr.conf.Signing.Method.Alg(),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(entities.Claims)
		},
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

	files.POST("/upload", ctr.UploadHandler)
	files.GET("", ctr.GetUserFiles)
	files.DELETE("/:id", ctr.DeleteFile)
}

func (ctr *Controller) UploadHandler(c echo.Context) error {
	token, ok := c.Get("access_token").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "unauthorized",
		})
	}

	claims, ok := token.Claims.(*entities.Claims)
	if !ok {
		return e.InternalServerError(errors.New("token claims are invalid"))
	}

	userIdStr := claims.Subject
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return e.Unauthorized("unauthorized")
	}

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

	uploadUrl, _, err := ctr.saveUC.Execute(ctx, fileHeader, userId)
	if err != nil {
		return err
	}

	response := map[string]string{
		"upload_url": uploadUrl,
	}

	return c.JSON(http.StatusOK, response)
}

func (ctr *Controller) GetUserFiles(c echo.Context) error {
	token, ok := c.Get("access_token").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "unauthorized",
		})
	}

	claims, ok := token.Claims.(*entities.Claims)
	if !ok {
		return e.InternalServerError(errors.New("token claims are invalid"))
	}

	userIdStr := claims.Subject
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return e.Unauthorized("unauthorized")
	}

	ctx := c.Request().Context()

	files, err := ctr.getFilesUC.Execute(ctx, userId)
	if err != nil {
		return err
	}

	response := map[string]interface{}{
		"data": files,
	}

	return c.JSON(http.StatusOK, response)
}

func (ctr *Controller) DeleteFile(c echo.Context) error {
	token, ok := c.Get("access_token").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{
			"error": "unauthorized",
		})
	}

	claims, ok := token.Claims.(*entities.Claims)
	if !ok {
		return e.InternalServerError(errors.New("token claims are invalid"))
	}

	userIdStr := claims.Subject
	_, err := uuid.Parse(userIdStr)
	if err != nil {
		return e.Unauthorized("unauthorized")
	}

	idStr := c.Param("id")
	if idStr == "" {
		return e.BadRequest("empty id")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return e.BadRequest("invalid id")
	}

	ctx := c.Request().Context()

	if err := ctr.deleteFileUC.Execute(ctx, id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
