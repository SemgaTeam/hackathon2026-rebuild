package errors

var (
	ErrPlaylistNotFound   = NewError("playlist not found")
	ErrInvalidDeleteRange = NewError("invalid delete range")
	ErrInvalidMoveRange   = NewError("invalid move range")
)
