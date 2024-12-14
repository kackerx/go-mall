package middleware

import (
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"

	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/common/util"
)

// 定义上下文键
const (
	TraceIDKey = "trace_id"
	SpanIDKey  = "span_id"
	TracerName = "gmall-service"
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

// TracingMiddleware Gin的链路追踪中间件
func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取tracer
		tracer := otel.Tracer(TracerName)

		// 创建span
		ctx, span := tracer.Start(c.Request.Context(), c.Request.URL.Path)
		defer span.End()

		// 获取trace信息
		traceID := span.SpanContext().TraceID().String()
		spanID := span.SpanContext().SpanID().String()

		// 将trace信息保存到gin上下文
		c.Set(TraceIDKey, traceID)
		c.Set(SpanIDKey, spanID)

		// 将context保存到gin上下文
		c.Request = c.Request.WithContext(ctx)

		// 处理请求
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
	logger.New(c).Info("AccessLog",
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

// GinPanicRecovery 自定义gin recover输出
func GinPanicRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.New(c).Error("http request broken pipe", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest))
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				logger.New(c).Error("http_request_panic", "path", c.Request.URL.Path, "error", err, "request", string(httpRequest), "stack", string(debug.Stack()))

				c.AbortWithError(http.StatusInternalServerError, err.(error))
			}
		}()
		c.Next()
	}
}
