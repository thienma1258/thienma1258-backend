package internal_errors

const DUPLCIATE_CODE = 1

var CODE_MESSAGE map[int]string = map[int]string{
	1: "name is duplicate",
}
var ERROR_DUPLICATE = NewInternalError(DUPLCIATE_CODE, CODE_MESSAGE[DUPLCIATE_CODE])

type InternalError struct {
	Code    int
	Message string
	Detail  string
}

func (err *InternalError) Error() string {
	return err.Message
}

func NewInternalError(code int, message string) *InternalError {
	return &InternalError{
		Code:    code,
		Message: message,
	}
}
