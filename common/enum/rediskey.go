package enum

// 项目名:模块名:键名
const (
	RedisKeyDemoOrderDetail = "gomall:demo:order_detail_%s"
)

const (
	RedisKeyAccessToken        = "gomall:user:access_token_%s"
	RedisKeyRefreshToken       = "gomall:user:refresh_token_%s"
	RedisKeyUserSession        = "gomall:user:session_%d"
	RediskeyTokenRefreshLock   = "gomall:user:token_refresh_lock_%s"
	RediskeyPasswordresetToken = "gomall:user:password_reset_token_%s"
)
