package request

type UserRegisterReq struct {
	UserName        string `json:"user_name,omitempty" binding:"required,email"`
	Password        string `json:"password,omitempty" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm,omitempty" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname,omitempty" binding:"max=30"`
	Slogan          string `json:"slogan,omitempty" binding:"max=30"`
	Avatar          string `json:"avatar,omitempty" binding:"max=100"`
}

type UserLoginReq struct {
	Body struct {
		UserName string `json:"user_name" binding:"required,e164|email"`
		Password string `json:"password" binding:"required,min=8"`
	}

	Header struct {
		Platform string `json:"platform" binding:"required,oneof=H5 APP"`
	}
}
