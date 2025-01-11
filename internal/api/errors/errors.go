package errors

type ErrorType string

const (
	ValidationError   ErrorType = "VALIDATION_ERROR"
	NotFoundError     ErrorType = "NOT_FOUND"
	UnauthorizedError ErrorType = "UNAUTHORIZED"
	InternalError     ErrorType = "INTERNAL_ERROR"
)

type APIError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
}

func NewAPIError(errorType ErrorType, message string) *APIError {
	return &APIError{
		Type:    errorType,
		Message: message,
	}
}
