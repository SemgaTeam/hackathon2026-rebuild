package errors

var (
	ErrInvalidUUID = NewError("invalid uuid")
	ErrInvalidName = NewError("invalid name")
)
