package util

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"time"
)

const aesKEY = "abcdefghijklmnlp"
const md5Len = 4  // MD5 的部分保留的字节数
const aesLen = 16 // aes 加密后的字节数，12-->16
// 将userId和MD5 揉到一起
// 类似于md5(userId+time)(4字节)+aes(userId+time)(16字节)，最终40个字符
func genAccessToken(uid int64) (string, error) {
	byteInfo := make([]byte, 12)
	binary.BigEndian.PutUint64(byteInfo, uint64(uid))
	binary.BigEndian.PutUint32(byteInfo[8:], uint32(time.Now().UnixNano()))
	encodeByte, err := AesEncrypt(byteInfo, []byte(aesKEY))
	if err != nil {
		return "", err
	}
	md5Byte := md5.Sum(byteInfo)
	data := append(md5Byte[0:md5Len], encodeByte...)
	return hex.EncodeToString(data), nil
}

func genRefreshToken(userId int64) (string, error) {
	return genAccessToken(userId)
}

func GenUserAuthToken(uid int64) (accessToken, refreshToken string, err error) {
	accessToken, err = genAccessToken(uid)
	if err != nil {
		return
	}
	refreshToken, err = genRefreshToken(uid)
	if err != nil {
		return
	}

	return
}

func GenPasswordResetToken(userId int64) (string, error) {
	// 与AccessToken使用同一规则, 必要时可以反解出userId
	return genAccessToken(userId)
}

func GenSessionId(userId int64) string {
	return fmt.Sprintf("%d-%d-%s", userId, time.Now().Unix(), RandNumStr(6))
}

// ParseUserIdFromToken 从Token中反解出userId,
// 后端服务redis不可用也没法立即恢复时可以使用这个方式保持产品最基本功能的使用, 不至于直接白屏
func ParseUserIdFromToken(accessToken string) (userId int64, err error) {
	if len(accessToken) != 2*(md5Len+aesLen) {
		// Token 格式不对
		return
	}
	encodeStr := accessToken[md5Len*2:]
	data, err := hex.DecodeString(encodeStr)
	if err != nil {
		return
	}
	decodeByte, _ := AesDecrypt(data, []byte(aesKEY)) // 忽略错误
	uid := binary.BigEndian.Uint64(decodeByte)
	if uid == 0 {
		return
	}
	userId = int64(uid)
	return
}
