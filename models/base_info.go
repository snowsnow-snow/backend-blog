package models

type BaseInfo struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	PublishIp   string    `json:"publishIp"`
	RequestId   string    `json:"requestId"`
	CreatedTime TimeStamp `json:"createdTime" gorm:"type:INTEGER"`
	UpdatedTime TimeStamp `json:"updatedTime" gorm:"type:INTEGER"`
	CreateUser  string    `json:"createUser"`
	UpdateUser  string    `json:"updateUser"`
}
