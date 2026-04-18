package errors

var (
	ErrPlaylistNotFound = NewError("playlist not found")
	ErrInvalidMoveRange = NewError("invalid move range")
)
