package errors

var (
	ErrInvalidExtension = NewError("invalid extension")
	ErrFileTooBig       = NewError("file too big")
)
