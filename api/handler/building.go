package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/request"
	"github.com/kackerx/go-mall/common/app"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/library"
	"github.com/kackerx/go-mall/logic/appservice"
	"github.com/kackerx/go-mall/logic/domainservice"
)

func TestErr(c *gin.Context) {
	// 使用Wrap包装错误, 生成项目错误
	err := errors.New("dao error")
	appErr := errcode.Wrap("包装错误", err)
	bAppErr := errcode.Wrap("再包装错误", appErr)
	logger.New(c).Error("记录错误", "err", bAppErr)

	// 预定义的错误, 追加错误根因
	err = errors.New("domain err")
	apiErr := errcode.ErrServer.WithCause(err)
	logger.New(c).Error("API错误", "err", apiErr)
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

func TestCreateDemoOrder(c *gin.Context) {
	req := new(request.DemoOrderCreateReq)
	if err := c.ShouldBind(req); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	req.UserID = 111
	demoDao := dao.NewDemoDao(c)
	domainSvc := domainservice.NewDemoDomainSvc(c, demoDao)
	svc := appservice.NewDemoAppSvc(c, domainSvc)

	order, err := svc.CreateDemoOrder(req)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).Success(order)
}

func TestWhoisLibReq(c *gin.Context) {
	detail, err := library.NewWhoisLib().GetHostIPDetail(c)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).Success(detail)
}

func TestMakeToken(c *gin.Context) {
	userDomainSvc := domainservice.NewUserDomainSvc()
	userAppSvc := appservice.NewUserAppSvc(userDomainSvc)

	resp, err := userAppSvc.GenToken(c)
	if err != nil {
		if errors.Is(err, errcode.ErrUserInvalid) {
			logger.New(c).Error("invalid user", "err", err)
			app.NewResponse(c).Error(errcode.ErrUserInvalid)
		} else {
			var appErr *errcode.AppError
			errors.As(err, &appErr)
			app.NewResponse(c).Error(appErr)
		}

		return
	}

	app.NewResponse(c).Success(resp)
}

func TestGetToken(c *gin.Context) {
	app.NewResponse(c).Success(map[string]any{
		"user_id":    c.GetInt64("userID"),
		"session_id": c.GetString("sessionID"),
	})
}

func TestRefreshToken(c *gin.Context) {
	refreshToken := c.Query("refresh_token")
	if refreshToken == "" {
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}

	svc := domainservice.NewUserDomainSvc()
	appSvc := appservice.NewUserAppSvc(svc)
	token, err := appSvc.RefreshToken(c, refreshToken)
	if err != nil {
		if errors.Is(err, errcode.ErrTooManyRequests) {
			app.NewResponse(c).Error(errcode.ErrTooManyRequests)
		} else {
			appErr := err.(*errcode.AppError)
			app.NewResponse(c).Error(appErr)
		}
		return
	}

	app.NewResponse(c).Success(token)
}
