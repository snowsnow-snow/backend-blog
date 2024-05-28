package models

type BaseInfo struct {
	ID          string   `gorm:"primaryKey" json:"id"`
	PublishIp   string   `json:"publishIp"`
	RequestId   string   `json:"requestId"`
	CreatedTime DateTime `json:"createdTime"`
	UpdatedTime DateTime `json:"updatedTime"`
	CreateUser  string   `json:"createUser"`
	UpdateUser  string   `json:"updateUser"`
}
