package reply

type TokenResp struct {
	AccessToken   string `json:"access_token,omitempty"`
	RefreshToken  string `json:"refresh_token,omitempty"`
	Duration      int64  `json:"duration,omitempty"`
	SrvCreateTime string `json:"srv_create_time,omitempty"`
}
