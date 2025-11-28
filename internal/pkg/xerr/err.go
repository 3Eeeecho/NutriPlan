package xerr

import "errors"

var (
	ErrUsernameExists          = errors.New("用户名已存在")
	ErrUserNotFound            = errors.New("用户不存在")
	ErrPasswordIncorrect       = errors.New("密码错误")
	ErrUserProfileIncomplete   = errors.New("用户档案不完整，请先完善健康档案")
	ErrInvalidToken            = errors.New("无效的token")
	ErrTokenGeneration         = errors.New("生成token失败")
	ErrTokenParsing            = errors.New("token解析失败")
	ErrUnexpectedSigningMethod = errors.New("意外的签名方法")
)
