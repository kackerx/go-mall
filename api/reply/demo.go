package reply

type DemoOrderResp struct {
	UserID    int64  `json:"user_id,omitempty"`
	Amount    string `json:"amount,omitempty"`
	Code      string `json:"code,omitempty"`
	State     int8   `json:"state,omitempty"`
	PaidAt    string `json:"paid_at,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}
