package do

import "time"

type DemoOrder struct {
	ID        int64     `json:"id,omitempty"`
	UserID    int64     `json:"user_id,omitempty"`
	Amount    string    `json:"amount"`
	Code      string    `json:"code"`
	State     int8      `json:"state"`
	IsDel     uint      `json:"is_del"`
	PaidAt    time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
