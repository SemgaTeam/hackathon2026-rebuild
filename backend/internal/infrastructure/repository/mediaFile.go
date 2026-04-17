package repository

import (
	"github.com/SemgaTeam/semga-stream/internal/config"
	"github.com/SemgaTeam/semga-stream/internal/core/entities"
	e "github.com/SemgaTeam/semga-stream/internal/core/errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"context"
	"errors"
)

type MediaFileRepository struct {
	conf *config.Config
	db   *pgxpool.Pool
}

func NewMediaFileRepository(conf *config.Config, pool *pgxpool.Pool) *MediaFileRepository {
	return &MediaFileRepository{
		conf,
		pool,
	}
}

func (r *MediaFileRepository) Save(ctx context.Context, media *entities.MediaFile) error {
	if media.ID != uuid.Nil {
		return r.Update(ctx, media)
	}

	return r.Create(ctx, media)
}

func (r *MediaFileRepository) Create(ctx context.Context, media *entities.MediaFile) error {
	var id uuid.UUID

	err := r.db.QueryRow(ctx,
		`INSERT INTO media_files(owner_id, type, file_name, file_path, file_size, mime_type, duration_seconds, created_at, status, is_deleted) 
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id`,
		media.OwnerID, media.Type, media.FileName, media.FilePath, media.FileSize, media.MimeType, media.DurationSeconds, media.CreatedAt, media.Status, media.IsDeleted,
	).Scan(&id)
	if err != nil {
		return e.Unknown(err)
	}

	media.ID = id

	return nil
}

func (r *MediaFileRepository) Update(ctx context.Context, media *entities.MediaFile) error {
	res, err := r.db.Exec(ctx,
		`UPDATE media_files
		 SET owner_id = $2,
				 type = $3,
				 file_name = $4,
				 file_path = $5,
				 file_size = $6,
				 mime_type = $7,
				 duration_seconds = $8,
				 created_at = $9,
				 status = $10,
				 is_deleted = $11
		 WHERE id = $1`,
		media.ID, media.OwnerID, media.Type, media.FileName, media.FilePath, media.FileSize, media.MimeType, media.DurationSeconds, media.CreatedAt, media.Status, media.IsDeleted,
	)
	if err != nil {
		return e.Unknown(err)
	}

	if res.RowsAffected() == 0 {
		return e.ErrFileNotFound
	}

	return nil
}

func (r *MediaFileRepository) ByID(ctx context.Context, id uuid.UUID) (*entities.MediaFile, error) {
	var mediaFile entities.MediaFile
	res := r.db.QueryRow(ctx,
		`SELECT id, owner_id, type, file_name, file_path, file_size, mime_type, duration_seconds, created_at, is_deleted, status
		 FROM media_files
		 WHERE id = $1`,
		id,
	)
	err := res.Scan(&mediaFile.ID, &mediaFile.OwnerID, &mediaFile.Type, &mediaFile.FileName, &mediaFile.FilePath, &mediaFile.FileSize, &mediaFile.MimeType, &mediaFile.DurationSeconds, &mediaFile.CreatedAt, &mediaFile.IsDeleted, &mediaFile.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, e.Unknown(err)
	}

	return &mediaFile, nil
}

func (r *MediaFileRepository) ByUserID(ctx context.Context, userID uuid.UUID) ([]entities.MediaFile, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, owner_id, type, file_name, file_path, file_size, mime_type, duration_seconds, created_at, is_deleted, status
		 FROM media_files
		 WHERE owner_id = $1`,
		userID,
	)
	if err != nil {
		return nil, e.Unknown(err)
	}
	defer rows.Close()

	var mediaFiles []entities.MediaFile

	for rows.Next() {
		var f entities.MediaFile
		err := rows.Scan(&f.ID, &f.OwnerID, &f.Type, &f.FileName, &f.FilePath, &f.FileSize, &f.MimeType, &f.DurationSeconds, &f.CreatedAt, &f.IsDeleted, &f.Status)
		if err != nil {
			return nil, e.Unknown(err)
		}
		mediaFiles = append(mediaFiles, f)
	}

	return mediaFiles, nil
}

func (r *MediaFileRepository) ByPath(ctx context.Context, path string) (*entities.MediaFile, error) {
	var mediaFile entities.MediaFile
	res := r.db.QueryRow(ctx,
		`SELECT id, owner_id, type, file_name, file_path, file_size, mime_type, duration_seconds, created_at, is_deleted, status
		 FROM media_files
		 WHERE file_path = $1`,
		path,
	)

	err := res.Scan(&mediaFile.ID, &mediaFile.OwnerID, &mediaFile.Type, &mediaFile.FileName, &mediaFile.FilePath, &mediaFile.FileSize, &mediaFile.MimeType, &mediaFile.DurationSeconds, &mediaFile.CreatedAt, &mediaFile.IsDeleted, &mediaFile.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, e.Unknown(err)
	}

	return &mediaFile, nil
}
