package response

type ErrorCode int

const (
	// 10xxx Account
	ErrCodeInvalidCredentials ErrorCode = 10001
	ErrCodeUserInactive       ErrorCode = 10002

	// 20xxx Conversation
	ErrCodeInvalidConvID        ErrorCode = 20001
	ErrCodeConversationNotFound ErrorCode = 20002

	// 30xxx Provider
	ErrCodeProviderNotFound ErrorCode = 30001

	// 40xxx Model
	ErrCodeModelNotFound ErrorCode = 40001

	// 90xxx 系统/通用
	ErrCodeInternalError  ErrorCode = 90000
	ErrCodeInvalidRequest ErrorCode = 90001
	ErrCodeUnauthorized   ErrorCode = 90002
)

func HTTPStatus(code ErrorCode) int {
	switch code / 1000 {
	case 10, 20, 30, 40:
		switch code {
		case ErrCodeInvalidCredentials:
			return 401
		case ErrCodeUserInactive:
			return 403
		case ErrCodeInvalidConvID:
			return 400
		case ErrCodeConversationNotFound, ErrCodeProviderNotFound, ErrCodeModelNotFound:
			return 404
		default:
			return 400
		}
	case 90:
		switch code {
		case ErrCodeInternalError:
			return 500
		case ErrCodeInvalidRequest:
			return 400
		case ErrCodeUnauthorized:
			return 401
		default:
			return 500
		}
	default:
		return 500
	}
}
