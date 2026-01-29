package dao

import (
	"backend-blog/internal/constant"
	"context"

	"gorm.io/gorm"
)

// GetDB 确保返回的 DB 实例绑定了当前的 Context
func GetDB(ctx context.Context) *gorm.DB {
	// 1. 优先尝试从 context 获取已有的事务
	if tx, ok := ctx.Value(constant.DBTxKey).(*gorm.DB); ok {
		// 即使是事务，也要 WithContext 以注入最新的 Context (包含 Trace ID 等)
		return tx.WithContext(ctx)
	}
	// 2. 如果没有事务，返回带 Context 的全局 DB
	return DB.WithContext(ctx)
}
