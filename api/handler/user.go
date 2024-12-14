package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/api/reply"
	"github.com/kackerx/go-mall/api/request"
	"github.com/kackerx/go-mall/common/app"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/logic/appservice"
	"github.com/kackerx/go-mall/logic/do"
)

type UserHandler struct {
	*Handler
	userAppSvc *appservice.UserAppSvc
}

func NewUserHandler(handler *Handler, userAppSvc *appservice.UserAppSvc) *UserHandler {
	return &UserHandler{Handler: handler, userAppSvc: userAppSvc}
}

func (uh *UserHandler) RegisterUser(c *gin.Context) {
	userRegisterReq := new(request.UserRegisterReq)
	if err := c.ShouldBind(userRegisterReq); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	if !util.PasswordComplexityVerify(userRegisterReq.Password) {
		logger.New(c).Warn("handler RegisterUser", "err", "密码复杂度不足", "password", userRegisterReq.Password)
		app.NewResponse(c).Error(errcode.ErrParams)
		return
	}

	userDo := new(do.UserBaseInfo)
	util.Copy(userDo, userRegisterReq)
	uid, err := uh.userAppSvc.UserRegister(c, userRegisterReq)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) {
			app.NewResponse(c).Error(errcode.ErrUserNameOccupied)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(&reply.UserRegisterResp{uid})
}

func (uh *UserHandler) LoginUser(c *gin.Context) {
	req := new(request.UserLoginReq)

	if err := c.ShouldBindJSON(&req.Body); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	if err := c.ShouldBindHeader(&req.Header); err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}

	resp, err := uh.userAppSvc.UserLogin(c, req)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNotRight) {
			app.NewResponse(c).Error(errcode.ErrUserNotRight)
		} else {
			app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		}
		return
	}

	app.NewResponse(c).Success(resp)
}

func (uh *UserHandler) LoginoutUser(c *gin.Context) {
	uid := c.GetInt64("user_id")
	platform := c.GetString("platform")

	if err := uh.userAppSvc.UserLoginout(c, uid, platform); err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}

	app.NewResponse(c).SuccessOK()
}
