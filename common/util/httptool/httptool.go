package httptool

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/util"
)

func Request(method, url string, options ...Option) (httpStatusCode int, resp []byte, err error) {
	start := time.Now()
	defer func() {
		fmt.Println(time.Since(start).Milliseconds())
	}()

	reqOption := defaultRequestOption()
	for _, opt := range options {
		opt(reqOption)
	}

	traceID, spanID, _ := util.GetTraceInfoFromCtx(reqOption.ctx)
	reqOption.headers["traceid"] = traceID
	reqOption.headers["spanid"] = spanID

	logger := logger.New(reqOption.ctx)
	defer func() {
		if err != nil {
			logger.Error("HTTP_REQUEST_ERROR_LOG", "method", method, "url", url, "body", reqOption.data, "reply", resp, "err", err)
		}
	}()

	req, err := http.NewRequest(method, url, bytes.NewReader(reqOption.data))
	if err != nil {
		return
	}

	reqOption.ctx, _ = context.WithTimeout(reqOption.ctx, reqOption.timeout)
	req.WithContext(reqOption.ctx)
	defer req.Body.Close()

	for k, v := range reqOption.headers {
		req.Header.Add(k, v)
	}

	// 发起请求
	client := http.Client{Timeout: reqOption.timeout}
	httpResp, err := client.Do(req)
	if err != nil {
		return
	}
	defer httpResp.Body.Close()

	// 记录请求日志
	dur := time.Since(start).Milliseconds()
	defer func() {
		if dur >= 3000 {
			logger.Warn("HTTP_REQUEST_SLOW_LOG", "method", method, "url", url, "body", reqOption.data, "reply", string(resp), "err", err, "dur/ms", dur, "header", req.Header)
		} else {
			logger.Debug("HTTP_REQUEST_DEBUG_LOG", "method", method, "url", url, "body", reqOption.data, "reply", string(resp), "err", err, "dur/ms", dur, "header", req.Header)
		}
	}()

	httpStatusCode = httpResp.StatusCode
	if httpStatusCode != http.StatusOK {
		err = errcode.Wrap("request api status not ok", errors.New(fmt.Sprintf("request status %d", httpStatusCode)))
		return
	}

	resp, _ = io.ReadAll(httpResp.Body)
	return
}

func Get(ctx context.Context, url string, options ...Option) (httpStatusCode int, resp []byte, err error) {
	options = append(options, WithContext(ctx))
	return Request(http.MethodGet, url, options...)
}

func Post(ctx context.Context, url string, data []byte, options ...Option) (httpStatusCode int, resp []byte, err error) {
	options = append(
		options,
		WithHeaders(map[string]string{"Content-Type": "application/json"}),
		WithContext(ctx),
		WithData(data),
	)

	return Request(http.MethodPost, url, options...)
}
