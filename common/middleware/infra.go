package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/log"
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

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存body
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))

		start := time.Now()
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw // 这里wrapper了一层, 让输出先写入到blw的body, 然后再让gin写入到自己的writer, 拿到响应
		accessLog(c, "access_start", time.Since(start), reqBody, nil)
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()
		c.Next()
		return
	}
}

func accessLog(c *gin.Context, accessType string, dur time.Duration, body []byte, out any) {
	req := c.Request
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path
	// todo: token记录
	log.New(c).Info("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", out,
		"time", int64(dur/time.Millisecond),
	)
}
