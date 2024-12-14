package domainservice

import (
	"context"
	"time"

	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/util"
	"github.com/kackerx/go-mall/dal/cache"
	"github.com/kackerx/go-mall/dal/dao"
	"github.com/kackerx/go-mall/logic/do"
)

type UserDomainSvc struct {
	userDao *dao.UserDao
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
	logger := logger.New(ctx)
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

	// 用旧refreshToken获取对应session, 如果是被窃取的refreshToken, 延迟过期也能获取到, 但是和最新的已经不同了
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

func (us *UserDomainSvc) RegisterUser(ctx context.Context, user *do.UserBaseInfo, password string) (userID int64, err error) {
	existUser, err := us.userDao.FindUserByUserName(ctx, user.UserName)
	if err != nil {
		return 0, errcode.Wrap("UserDomainSvc RegisterUser err", err)
	}

	if existUser.UserName != "" {
		return 0, errcode.ErrUserNameOccupied
	}

	bcryptPassword, err := util.BcryptPassword(password)
	if err != nil {
		return 0, errcode.Wrap("UserDomainSvc RegisterUser BcryptPassword err", err)
	}

	return us.userDao.CreateUser(ctx, user, bcryptPassword)
}

func (us *UserDomainSvc) LoginUser(ctx context.Context, userName, password, platform string) (tokenInfo *do.TokenInfo, err error) {
	user, err := us.userDao.FindUserByUserName(ctx, userName)
	if err != nil {
		return nil, errcode.Wrap("UserDomainSvc LoginUser err", err)
	}

	if user.UserName == "" {
		return nil, errcode.ErrUserNotRight
	}

	if !util.BcryptCompare(user.Password, password) {
		return nil, errcode.ErrUserNotRight
	}

	return us.GenAuthToken(ctx, int64(user.ID), platform, "")
}

func (us *UserDomainSvc) LoginoutUser(ctx context.Context, userID int64, platform string) (err error) {
	session, err := cache.GetUserPlatformSession(ctx, userID, platform)
	if err != nil {
		return errcode.Wrap("UserDomainSvc LoginoutUser GetUserPlatformSession err", err)
	}

	if err = cache.DelAccessToken(ctx, session.AccessToken); err != nil {
		return errcode.Wrap("UserDomainSvc LoginoutUser DelAccessToken err", err)
	}

	if err = cache.DelRefreshToken(ctx, session.RefreshToken); err != nil {
		return errcode.Wrap("UserDomainSvc LoginoutUser DelRefreshToken err", err)
	}

	if err = cache.DelUserSessionOnPlatform(ctx, userID, platform); err != nil {
		return errcode.Wrap("UserDomainSvc LoginoutUser DelUserSessionOnPlatform err", err)
	}

	return
}
