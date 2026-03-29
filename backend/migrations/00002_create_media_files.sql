-- +goose Up
-- +goose StatementBegin
CREATE TABLE media_files (
  id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type             VARCHAR(10) NOT NULL, -- audio / video
  file_name        VARCHAR(255) NOT NULL,
  file_path        TEXT NOT NULL,
  file_size        BIGINT NOT NULL,
  mime_type        VARCHAR(100) NOT NULL,
  duration_seconds INT,
  created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_deleted       BOOLEAN NOT NULL DEFAULT FALSE
);

-- Type limits
ALTER TABLE media_files
ADD CONSTRAINT chk_media_type
CHECK (type IN ('audio', 'video'));

-- Size limits
ALTER TABLE media_files
ADD CONSTRAINT chk_media_size
CHECK (
    (type = 'audio' AND file_size <= 52428800) -- 50MB
    OR
    (type = 'video' AND file_size <= 1048576000) -- 1000MB
);

-- Indexes
CREATE INDEX idx_media_owner ON media_files(owner_id);
CREATE INDEX idx_media_type ON media_files(type);
CREATE INDEX idx_media_created_at ON media_files(created_at);
CREATE INDEX idx_media_is_deleted ON media_files(is_deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS media_files CASCADE;
-- +goose StatementEnd
