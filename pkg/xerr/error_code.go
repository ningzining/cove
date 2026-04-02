package xerr

/**(前3位代表业务,后三位代表具体功能)**/
var (
	ErrCommon       = 100000 // 通用错误码
	ErrBind         = 100001 // 绑定参数错误码
	ErrInvalidParam = 100002 // 无效参数错误码
	ErrDB           = 100003 // 数据库错误码
	ErrTokenSign    = 100004 // token签名错误码
	ErrTokenInvalid = 100005 // token无效错误码
	ErrTokenExpired = 100006 // token过期错误码

	ErrPhoneExist   = 101000 // 手机号已存在错误码
	ErrLoginFailed  = 101001 // 登录失败错误码
	ErrUserDisabled = 101002 // 用户已禁用错误码
)
