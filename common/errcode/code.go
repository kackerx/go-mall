package errcode

// 预定义的错误是最终控制器层返回给请求的客户端的, 会封装统一的响应组件来处理
var (
	Success            = newError(0, "success")
	ErrServer          = newError(10000000, "服务器内部错误")
	ErrParams          = newError(10000001, "参数错误, 请检查")
	ErrNotFound        = newError(10000002, "资源未找到")
	ErrPanic           = newError(10000003, "系统开小差了~ 请稍后重试")
	ErrToken           = newError(10000004, "权限校验失败")
	ErrForbidden       = newError(10000005, "未授权")
	ErrTooManyRequests = newError(10000006, "请求过多")
)
