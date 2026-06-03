package response

import "omni-pixel/domain"

var (
	ErrInvalidCredentials    = New(ErrCodeInvalidCredentials, "邮箱或密码错误")
	ErrUserInactive          = New(ErrCodeUserInactive, "账户未激活")
	ErrUserAlreadyExists     = New(ErrCodeUserAlreadyExists, "用户已存在")
	ErrInvalidConvID         = New(ErrCodeInvalidConvID, "会话 ID 格式错误")
	ErrConversationNotFound  = New(ErrCodeConversationNotFound, "会话不存在")
	ErrProviderNotFound      = New(ErrCodeProviderNotFound, "厂商不存在")
	ErrModelNotFound         = New(ErrCodeModelNotFound, "模型不存在")
	ErrInternalServer        = New(ErrCodeInternalError, "服务器内部错误")
	ErrInvalidRequest        = New(ErrCodeInvalidRequest, "请求参数无效")
	ErrUnauthorized          = New(ErrCodeUnauthorized, "认证失败")
)

var domainMappings = map[error]*APIError{
	domain.ErrInvalidCredentials:    ErrInvalidCredentials,
	domain.ErrUserInactive:          ErrUserInactive,
	domain.ErrUserAlreadyExists:     ErrUserAlreadyExists,
	domain.ErrConversationNotFound:  ErrConversationNotFound,
	domain.ErrInvalidConversationID: ErrInvalidConvID,
	domain.ErrProviderNotFound:      ErrProviderNotFound,
	domain.ErrModelNotFound:         ErrModelNotFound,
}
