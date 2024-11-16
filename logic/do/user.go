package do

import "time"

type SessionInfo struct {
	UserID       int64  `json:"user_id,omitempty"`
	Platform     string `json:"platform,omitempty"`
	SessionID    string `json:"session_id,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenInfo struct {
	AccessToken   string    `json:"access_token,omitempty"`
	RefreshToken  string    `json:"refresh_token,omitempty"`
	Duration      int64     `json:"duration,omitempty"`
	SrvCreateTime time.Time `json:"srv_create_time"`
}

type UserBaseInfo struct {
	ID        uint64    `json:"id,omitempty"`
	Nickname  string    `json:"nickname,omitempty"`
	UserName  string    `json:"user_name,omitempty"`
	Verified  uint      `json:"verified,omitempty"`
	Avatar    string    `json:"avatar,omitempty"`
	Slogan    string    `json:"slogan,omitempty"`
	IsBlocked uint      `json:"is_blocked,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenVerify struct {
	Approved  bool   `json:"approved,omitempty"`
	UserID    int64  `json:"user_id,omitempty"`
	SessionID string `json:"session_id,omitempty"`
}
