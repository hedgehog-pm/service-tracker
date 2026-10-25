package domain

type Code string

const (
	CodeInvalidInput  Code = "invalid_input"
	CodeNotFound      Code = "not_found"
	CodeConflict      Code = "conflict"
	CodeInternalError Code = "internal_error"
)

type Error struct {
	Code    Code
	Message string
}

func NewError(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrProjectNotFound     = NewError(CodeNotFound, "project not found")
	ErrTaskNotFound        = NewError(CodeNotFound, "task not found")
	ErrStatusNotFound      = NewError(CodeNotFound, "status not found")
	ErrModuleNotFound      = NewError(CodeNotFound, "module not found")
	ErrInvalidReference    = NewError(CodeInvalidInput, "referenced resource is missing or belongs to another project")
	ErrResourceInUse       = NewError(CodeConflict, "resource is in use")
	ErrAlreadyExists       = NewError(CodeConflict, "resource already exists")
	ErrInternalServerError = NewError(CodeInternalError, "internal error")
)
