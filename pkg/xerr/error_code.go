package xerr

/**(前3位代表业务,后三位代表具体功能)**/
var (
	ErrCommon       = 100000 // 通用错误
	ErrBind         = 100001 // 绑定参数错误
	ErrInvalidParam = 100002 // 无效参数错误
	ErrDB           = 100003 // 数据库错误
	ErrTokenSign    = 100004 // token签名错误
	ErrTokenInvalid = 100005 // token无效
	ErrTokenExpired = 100006 // token过期
	ErrUnauthorized = 100007 // 未认证
	ErrForbidden    = 100008 // 未授权

	ErrLoginFailed     = 101000 // 登录失败
	ErrAccountDisabled = 101001 // 账号已禁用

	ErrRoleCodeExist     = 102000 // 角色编码已存在
	ErrRoleNameExist     = 102001 // 角色名称已存在
	ErrRoleNotExist      = 102002 // 角色不存在
	ErrRoleInUse         = 102003 // 角色已被用户使用，无法删除
	ErrRoleCannotOperate = 102004 // 内置角色，不能操作
)
