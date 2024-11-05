package do

// UserRegister 用户注册信息
type UserRegister struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
}
