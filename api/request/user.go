package request

type UserRegisterReq struct {
	UserName        string `json:"user_name,omitempty" binding:"required,email"`
	Password        string `json:"password,omitempty" binding:"required,min=8"`
	PasswordConfirm string `json:"password_confirm,omitempty" binding:"required,eqfield=Password"`
	Nickname        string `json:"nickname,omitempty" binding:"max=30"`
	Slogan          string `json:"slogan,omitempty" binding:"max=30"`
	Avatar          string `json:"avatar,omitempty" binding:"max=100"`
}
