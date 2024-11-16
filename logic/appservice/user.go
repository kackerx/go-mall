package appservice

import (
	"context"
	"errors"

	"github.com/kackerx/go-mall/api/reply"
	"github.com/kackerx/go-mall/api/request"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/logic/do"
	"github.com/kackerx/go-mall/logic/domainservice"
)

type UserAppSvc struct {
	userDomainSvc *domainservice.UserDomainSvc
}

func NewUserAppSvc(userDomainSvc *domainservice.UserDomainSvc) *UserAppSvc {
	return &UserAppSvc{userDomainSvc: userDomainSvc}
}
func (us *UserAppSvc) GenToken(ctx context.Context) (resp *reply.TokenResp, err error) {
	tokenInfo, err := us.userDomainSvc.GenAuthToken(ctx, 110, "h5", "")
	if err != nil {
		return nil, err
	}

	log.New(ctx).Info("gen token success", tokenInfo)

	resp = new(reply.TokenResp)
	if err = util.Copy(resp, tokenInfo); err != nil {
		return
	}

	return
}

func (us *UserAppSvc) RefreshToken(ctx context.Context, refreshToken string) (resp *reply.TokenResp, err error) {
	token, err := us.userDomainSvc.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	log.New(ctx).Info("token refresh success", "token", token)
	resp = new(reply.TokenResp)
	if err := util.Copy(resp, token); err != nil {
		return nil, err
	}

	return
}

func (us *UserAppSvc) UserRegister(ctx context.Context, req *request.UserRegisterReq) (int64, error) {
	user := new(do.UserBaseInfo)
	util.Copy(user, req)

	id, err := us.userDomainSvc.RegisterUser(ctx, user, req.Password)
	if err != nil {
		if errors.Is(err, errcode.ErrUserNameOccupied) { // 用户名已存在不用额外的处理
			return 0, err
		}

		// todo: 外围逻辑通知用户注册失败, 日志, 告警, 提示灯
		return 0, err
	}

	// todo: 注册成功后的发消息通知, 事件推送, 或者直接跳转登录逻辑, 非核心的应用层的职责

	return id, nil
}

func (us *UserAppSvc) UserLogin(ctx context.Context, req *request.UserLoginReq) (resp *reply.TokenResp, err error) {
	tokenInfo, err := us.userDomainSvc.LoginUser(ctx, req.Body.UserName, req.Body.Password, req.Header.Platform)
	if err != nil {
		return
	}

	resp = new(reply.TokenResp)
	util.Copy(resp, tokenInfo)
	// todo: 登录成功的外围逻辑
	return
}

func (us *UserAppSvc) UserLoginout(ctx context.Context, userID int64, platform string) (err error) {
	return us.userDomainSvc.LoginoutUser(ctx, userID, platform)
}
