package domainservice

import (
	"context"
	"time"

	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/log"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/cache"
	"github.com/kackerx/go-mall/logic/do"
)

type UserDomainSvc struct {
}

func NewUserDomainSvc() *UserDomainSvc {
	return &UserDomainSvc{}
}

func (us *UserDomainSvc) GetUserBaseInfo(userID int64) *do.UserBaseInfo {
	return &do.UserBaseInfo{
		ID:        110,
		Nickname:  "kacker",
		UserName:  "kingvstr@hotmail.com",
		Verified:  1,
		Avatar:    "",
		Slogan:    "",
		IsBlocked: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (us *UserDomainSvc) GenAuthToken(ctx context.Context, userID int64, platform, sessionID string) (tokenInfo *do.TokenInfo, err error) {
	user := us.GetUserBaseInfo(userID)

	// 异常case
	if user.ID == 0 || user.IsBlocked == enum.UserBlockStateBlocked {
		err = errcode.ErrUserInvalid
		return
	}

	userSession := new(do.SessionInfo)
	userSession.UserID = userID
	userSession.Platform = platform
	if sessionID == "" {
		// 为空是登录行为, 生成会话id
		sessionID = util.GenSessionId(userID)
	}
	userSession.SessionID = sessionID
	accessToken, refreshToken, err := util.GenUserAuthToken(userID)
	if err != nil {
		err = errcode.Wrap("Token生成失败", err)
		return
	}

	userSession.AccessToken = accessToken
	userSession.RefreshToken = refreshToken

	// 删除旧token
	if err = cache.DelOldSessionTokens(ctx, userSession); err != nil {
		err = errcode.Wrap("DelOldSessionTokens err", err)
		return
	}

	// 设置缓存
	if err = cache.SetUserToken(ctx, userSession); err != nil {
		err = errcode.Wrap("设置token缓存错误", err)
		return
	}

	if err = cache.SetUserSession(ctx, userSession); err != nil {
		err = errcode.Wrap("设置Session错误", err)
		return
	}

	return &do.TokenInfo{
		AccessToken:   accessToken,
		RefreshToken:  refreshToken,
		Duration:      int64(enum.AccessTokenDuration.Seconds()),
		SrvCreateTime: time.Now(),
	}, nil
}

func (us *UserDomainSvc) RefreshToken(ctx context.Context, refreshToken string) (resp *do.TokenInfo, err error) {
	logger := log.New(ctx)
	ok, err := cache.LockTokenRefresh(ctx, refreshToken)
	defer cache.UnlockTokenRefresh(ctx, refreshToken)

	if err != nil {
		err = errcode.Wrap("设置refresh锁失败", err)
		return
	}

	if !ok {
		err = errcode.ErrTooManyRequests
		return
	}

	session, err := cache.GetRefreshToken(ctx, refreshToken)
	if err != nil || session == nil {
		logger.Error("GetRefreshToken faild", "err", err)
		err = errcode.ErrToken
		return
	}

	platformSession, err := cache.GetUserPlatformSession(ctx, session.UserID, session.Platform)
	if err != nil {
		logger.Error("GetUserPlatformSession faild", "err", err)
		err = errcode.ErrToken
		return
	}

	// 请求刷新的re和session中的不一致, 说明re过时, 可能是re被窃取或者前端刷token并发控制问题
	if platformSession.RefreshToken != refreshToken {
		logger.Warn("ExpiredRefreshToken", "request_token", refreshToken, "new_token", platformSession.RefreshToken, "user_id", platformSession.UserID)
		err = errcode.ErrToken
		return
	}

	// 重新生成token
	tokenInfo, err := us.GenAuthToken(ctx, session.UserID, session.Platform, session.SessionID)
	if err != nil {
		err = errcode.Wrap("GenUserAuthToken err", err)
		return
	}

	return tokenInfo, nil
}
