-- +goose Up
-- +goose StatementBegin
CREATE TABLE playlists (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id   UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name       VARCHAR(50) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE playlist_items (
  playlist_id   UUID NOT NULL REFERENCES playlists(id),
  media_file_id UUID NOT NULL REFERENCES media_files(id),
  position      INTEGER NOT NULL,
  PRIMARY KEY (playlist_id, media_file_id),
  UNIQUE (playlist_id, position)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS playlist_items CASCADE;
DROP TABLE IF EXISTS playlists CASCADE;
-- +goose StatementEnd
