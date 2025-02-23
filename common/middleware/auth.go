package middleware

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/kackerx/go-mall/common/app"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/dal/cache"
	"github.com/kackerx/go-mall/logic/do"
)

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("gomall-token")
		if len(token) != 40 {
			app.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}

		tokenVerify, err := VerifyAccessToken(c, token)
		if err != nil || !tokenVerify.Approved {
			app.NewResponse(c).Error(errcode.ErrToken)
			c.Abort()
			return
		}

		c.Set("user_id", tokenVerify.UserID)
		c.Set("session_id", tokenVerify.SessionID)
		c.Set("platform", tokenVerify.Platform)
		c.Next()
	}
}

// VerifyAccessToken 校验token合法
func VerifyAccessToken(ctx context.Context, accessToken string) (resp *do.TokenVerify, err error) {
	tokenInfo, err := cache.GetAccessToken(ctx, accessToken)
	if err != nil {
		logger.New(ctx).Error("GetAccessToken err", "err", err)
		return nil, err
	}

	resp = new(do.TokenVerify)
	if tokenInfo != nil && tokenInfo.UserID != 0 {
		resp.Approved = true
		resp.UserID = tokenInfo.UserID
		resp.SessionID = tokenInfo.SessionID
		resp.Platform = tokenInfo.Platform
	}

	return
}
