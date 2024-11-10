package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/util"
)

func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("traceid")
		pSpanID := c.Request.Header.Get("spanid")
		// 网关在调用其他服务完成业务逻辑时, 生成自己的spanID, 带上上一个服务的spanid作为pspanid
		spanID := util.GenerateSpanID(c.Request.RemoteAddr)

		if traceID == "" { // traceid为空证明是链路的起始, 设置为此次的spanID
			traceID = spanID // trace标识整个请求链路, span表示链路中的不同服务
		}

		c.Set("traceid", traceID)
		c.Set("spanid", spanID)
		c.Set("pspanid", pSpanID)
		c.Next()
	}
}
