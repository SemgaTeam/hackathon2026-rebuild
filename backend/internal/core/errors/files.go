package errors

var (
	ErrFileNotFound     = NewError("file not found")
	ErrOpeningFile      = NewError("error opening file")
	ErrInvalidExtension = NewError("invalid extension")
	ErrInvalidMIMEType  = NewError("invalid mime type")
	ErrFileTooBig       = NewError("file too big")
)
