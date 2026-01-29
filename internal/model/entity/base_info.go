package entity

import (
	"backend-blog/internal/model"
	"backend-blog/internal/pkg"
	"backend-blog/utility"

	"gorm.io/gorm"
)

type BaseInfo struct {
	ID          int64           `Gorm:"primaryKey;autoIncrement:false" json:"id,string"` // Snowflake ID (Json String安全)
	PublishIp   string          `json:"publishIp"`
	TraceId     string          `json:"traceId"`
	CreatedTime model.TimeStamp `json:"createdTime" Gorm:"type:INTEGER"`
	UpdatedTime model.TimeStamp `json:"updatedTime" Gorm:"type:INTEGER"`
	CreateUser  string          `json:"createUser"`
	UpdateUser  string          `json:"updateUser"`
}

// BeforeCreate GORM Hook to generate Snowflake ID and populate metadata
func (b *BaseInfo) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == 0 {
		b.ID = utility.GenID()
	}
	if b.CreatedTime.IsZero() {
		b.CreatedTime = model.Now()
	}
	if b.UpdatedTime.IsZero() {
		b.UpdatedTime = model.Now()
	}

	// 提取 Context 信息
	ctx := tx.Statement.Context
	if b.TraceId == "" {
		if tid, ok := ctx.Value(pkg.TraceKey).(string); ok {
			b.TraceId = tid
		}
	}
	if b.PublishIp == "" {
		if ip, ok := ctx.Value(pkg.IPKey).(string); ok {
			b.PublishIp = ip
		}
	}
	return
}

// BeforeUpdate GORM Hook to populate metadata
func (b *BaseInfo) BeforeUpdate(tx *gorm.DB) (err error) {
	if b.UpdatedTime.IsZero() {
		b.UpdatedTime = model.Now()
	}

	// 提取 Context 信息
	ctx := tx.Statement.Context
	if tid, ok := ctx.Value(pkg.TraceKey).(string); ok {
		b.TraceId = tid
	}
	if ip, ok := ctx.Value(pkg.IPKey).(string); ok {
		b.PublishIp = ip
	}
	return
}
