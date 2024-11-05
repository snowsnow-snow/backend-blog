package entity

// User 用户表信息
type User struct {
	BaseInfo
	Username string
	Password string
	Salt     string
	Birthday string
	Avatar   string
	//Phone     int64  `gorm:"DEFAULT:0"`
	//Email     string `gorm:"type:varchar(20);unique_index;"`
}

// LoginRequest 用户登录请求信息
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
