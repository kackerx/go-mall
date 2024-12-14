package app

import (
	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code,omitempty"`
	Msg        string      `json:"msg,omitempty"`
	RequestID  string      `json:"request_id,omitempty"`
	Data       any         `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

func (r *response) Success(data any) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	if v, ok := r.ctx.Get("traceid"); ok {
		r.RequestID = v.(string)
	}

	r.Data = data
	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOK() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	if v, ok := r.ctx.Get("traceid"); ok {
		r.RequestID = v.(string)
	}

	// 兜底错误响应日志
	logger.New(r.ctx).Error("api_response_error", "error", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}

func (r *response) SetPagination(pagination *Pagination) *response {
	r.Pagination = pagination
	return r
}

func NewResponse(c *gin.Context) *response {
	return &response{ctx: c}
}
