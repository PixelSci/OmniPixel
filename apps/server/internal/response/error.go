package response

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func New(code ErrorCode, msg string) *APIError {
	return &APIError{Code: code, Message: msg}
}

func (e *APIError) Error() string {
	return e.Message
}
