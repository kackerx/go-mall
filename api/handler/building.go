package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/app"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/log"
)

func TestErr(c *gin.Context) {
	// 使用Wrap包装错误, 生成项目错误
	err := errors.New("dao error")
	appErr := errcode.Wrap("包装错误", err)
	bAppErr := errcode.Wrap("再包装错误", appErr)
	log.New(c).Error("记录错误", "err", bAppErr)

	// 预定义的错误, 追加错误根因
	err = errors.New("domain err")
	apiErr := errcode.ErrServer.WithCause(err)
	log.New(c).Error("API错误", "err", apiErr)
	c.JSON(apiErr.HttpStatusCode(), gin.H{
		"code": apiErr.Code(),
		"msg":  apiErr.Msg(),
	})
}

func TestRespErr(c *gin.Context) {
	baseErr := errors.New("dao err")
	err := errcode.Wrap("getUserService error", baseErr)
	app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
}

func TestRespSuccess(c *gin.Context) {
	data := []struct {
		Name string
		Age  int
	}{
		{
			Name: "kacker",
			Age:  28,
		},
	}

	pagination := app.NewPagination(c).SetTotal(100)
	app.NewResponse(c).
		SetPagination(pagination).
		Success(data)
}
