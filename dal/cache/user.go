package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/kackerx/go-mall/common/enum"
	"github.com/kackerx/go-mall/common/errcode"
	"github.com/kackerx/go-mall/common/logger"
	"github.com/kackerx/go-mall/logic/do"
)

// SetUserToken 设置用户的AccessToken 和 RefreshToken 缓存
func SetUserToken(ctx context.Context, session *do.SessionInfo) error {
	logger := logger.New(ctx)
	err := setAccessToken(ctx, session)
	if err != nil {
		logger.Error("redis error", "err", err)
		return err
	}
	err = setRefreshToken(ctx, session)
	if err != nil {
		logger.Error("redis error", "err", err)
		return err
	}
	return err
}

func SetUserSession(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.RedisKeyUserSession, session.UserID)
	sessionDataBytes, _ := json.Marshal(session)
	err := Redis().HSet(ctx, redisKey, session.Platform, sessionDataBytes).Err()
	if err != nil {
		logger.New(ctx).Error("redis error", "err", err)
		return err
	}
	return err
}

// DelOldSessionTokens 删除用户旧Session的Token
func DelOldSessionTokens(ctx context.Context, session *do.SessionInfo) error {
	// log := logger.New(ctx)
	oldSession, err := GetUserPlatformSession(ctx, session.UserID, session.Platform)
	if err != nil {
		return err
	}
	if oldSession == nil {
		// 没有旧Session
		return nil
	}
	err = DelAccessToken(ctx, oldSession.AccessToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	err = DelayDelRefreshToken(ctx, oldSession.RefreshToken)
	if err != nil {
		return errcode.Wrap("redis error", err)
	}
	return nil
}

// GetUserPlatformSession 获取用户在指定平台中的Session信息
func GetUserPlatformSession(ctx context.Context, userId int64, platform string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.RedisKeyUserSession, userId)
	result, err := Redis().HGet(ctx, redisKey, platform).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	// key 不存在
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	session := new(do.SessionInfo)
	err = json.Unmarshal([]byte(result), &session)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func setAccessToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.RedisKeyAccessToken, session.AccessToken)
	sessionDataBytes, _ := json.Marshal(session)
	res, err := Redis().Set(ctx, redisKey, sessionDataBytes, enum.AccessTokenDuration).Result()
	logger.New(ctx).Debug("redis debug", "res", res, "err", err)
	return err
}

func setRefreshToken(ctx context.Context, session *do.SessionInfo) error {
	redisKey := fmt.Sprintf(enum.RedisKeyRefreshToken, session.RefreshToken)
	sessionDataBytes, _ := json.Marshal(session)
	return Redis().Set(ctx, redisKey, sessionDataBytes, enum.RefreshTokenDuration).Err()
}

func DelAccessToken(ctx context.Context, accessToken string) error {
	redisKey := fmt.Sprintf(enum.RedisKeyAccessToken, accessToken)
	return Redis().Del(ctx, redisKey).Err()
}

// DelayDelRefreshToken 刷新Token时让旧的RefreshToken 保留一段时间自己过期
func DelayDelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.RedisKeyRefreshToken, refreshToken)
	return Redis().Expire(ctx, redisKey, enum.OldRefreshTokenHoldingDuration).Err()
}

// DelRefreshToken 直接删除RefreshToken缓存  修改密码、退出登录时使用
func DelRefreshToken(ctx context.Context, refreshToken string) error {
	redisKey := fmt.Sprintf(enum.RedisKeyRefreshToken, refreshToken)
	return Redis().Del(ctx, redisKey).Err()
}

// DelUserSessionOnPlatform Delete user's session on specific platform
func DelUserSessionOnPlatform(ctx context.Context, userId int64, platform string) error {
	redisKey := fmt.Sprintf(enum.RedisKeyUserSession, userId)
	return Redis().HDel(ctx, redisKey, platform).Err()
}

// DelUserSessions Delete user's sessions on all platform
func DelUserSessions(ctx context.Context, userId int64) error {
	// 先获取所有平台上的Session信息中
	sessions, err := GetUserAllSessions(ctx, userId)
	if err != nil {
		return err
	}
	// 把所有Session中保存的正在用的Token都过期掉
	for _, sessInfo := range sessions {
		DelOldSessionTokens(ctx, sessInfo)
	}
	// Token过期完成后再删掉Session
	redisKey := fmt.Sprintf(enum.RedisKeyUserSession, userId)
	return Redis().Del(ctx, redisKey).Err()
}

// GetUserAllSessions 获取用户在所有platform上的Session
func GetUserAllSessions(ctx context.Context, userId int64) (map[string]*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.RedisKeyUserSession, userId)
	result, err := Redis().HGetAll(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	// key 不存在
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	sessions := make(map[string]*do.SessionInfo)
	for platform, sessionData := range result {
		session := new(do.SessionInfo)
		err = json.Unmarshal([]byte(sessionData), &session)
		if err != nil {
			return nil, err
		}
		sessions[platform] = session
	}
	// logger.New(ctx).Debug("hgetall user all session", "data", sessions)
	return sessions, nil
}

func LockTokenRefresh(ctx context.Context, refreshToken string) (bool, error) {
	redisLockKey := fmt.Sprintf(enum.RediskeyTokenRefreshLock, refreshToken)
	return Redis().SetNX(ctx, redisLockKey, "locked", 10*time.Second).Result()
}

func UnlockTokenRefresh(ctx context.Context, refreshToken string) error {
	redisLockKey := fmt.Sprintf(enum.RediskeyTokenRefreshLock, refreshToken)
	return Redis().Del(ctx, redisLockKey).Err()
}

func GetRefreshToken(ctx context.Context, refreshToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.RedisKeyRefreshToken, refreshToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)

	return session, nil
}

func GetAccessToken(ctx context.Context, accessToken string) (*do.SessionInfo, error) {
	redisKey := fmt.Sprintf(enum.RedisKeyAccessToken, accessToken)
	result, err := Redis().Get(ctx, redisKey).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	session := new(do.SessionInfo)
	if errors.Is(err, redis.Nil) {
		return session, nil
	}
	json.Unmarshal([]byte(result), &session)

	return session, nil
}

// SetPasswordResetToken 设置重置密码的验证Token信息到缓存, 15分钟内有效
// @param ctx
// @param userId
// @param token 重置密码的验证Token
// @param code 验证码
func SetPasswordResetToken(ctx context.Context, userId int64, token, code string) error {
	redisKey := fmt.Sprintf(enum.RediskeyPasswordresetToken, token)
	val := fmt.Sprintf("%d:%s", userId, code) // val 以 userId:code 的字符串形式存储
	return Redis().Set(ctx, redisKey, val, enum.PasswordTokenDuration).Err()
}

func GetPasswordResetToken(ctx context.Context, token string) (userId int64, code string, err error) {
	redisKey := fmt.Sprintf(enum.RediskeyPasswordresetToken, token)
	val, redisErr := Redis().Get(ctx, redisKey).Result()
	if redisErr != nil && !errors.Is(redisErr, redis.Nil) {
		err = redisErr
		return
	}
	valArr := strings.Split(val, ":")
	if len(valArr) != 2 { // 密码重置Token无对应的缓存, 判定该参数不合法, 此处直接返回
		return
	}
	userId, _ = strconv.ParseInt(valArr[0], 10, 64)
	code = valArr[1]

	return
}

func DelPasswordResetToken(ctx context.Context, token string) error {
	redisKey := fmt.Sprintf(enum.RediskeyPasswordresetToken, token)
	return Redis().Del(ctx, redisKey).Err()
}
