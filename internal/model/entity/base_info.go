package entity

import (
	"backend-blog/internal/model"
)

type BaseInfo struct {
	ID          string          `gorm:"primaryKey" json:"id"`
	PublishIp   string          `json:"publishIp"`
	RequestId   string          `json:"requestId"`
	CreatedTime model.TimeStamp `json:"createdTime" gorm:"type:INTEGER"`
	UpdatedTime model.TimeStamp `json:"updatedTime" gorm:"type:INTEGER"`
	CreateUser  string          `json:"createUser"`
	UpdateUser  string          `json:"updateUser"`
}
