package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID        int64                 `json:"id,omitempty"`
	Nickname  string                `json:"nickname,omitempty"`
	UserName  string                `json:"user_name,omitempty"`
	Password  string                `json:"password,omitempty"`
	Verified  int                   `json:"verified,omitempty"`
	Avatar    string                `json:"avatar,omitempty"`
	Slogan    string                `json:"slogan,omitempty"`
	IsDel     soft_delete.DeletedAt `json:"is_del,omitempty"`
	IsBlocked int                   `json:"is_blocked,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

func (u *User) TableName() string {
	return "users"
}
