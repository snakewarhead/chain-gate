package errors

type AppError struct {
	Error   error
	Message string
	Code    int
	Custom  interface{}
}

func NewAppError(error error, message string, code int64, custom interface{}) *AppError {
	return &AppError{
		error,
		message,
		int(code),
		custom,
	}
}

func MarshallingError(err error) *AppError {
	return NewAppError(err, "", -1, nil)
}

func UnsupportedOperation() *AppError {
	return &AppError{nil, "unsupported operation", 14, nil}
}
