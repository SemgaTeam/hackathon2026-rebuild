package errors

var (
	ErrInvalidExtension = NewError("invalid extension")
	ErrInvalidMIMEType  = NewError("invalid mime type")
	ErrFileTooBig       = NewError("file too big")
)
