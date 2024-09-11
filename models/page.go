package models

type Page[T any] struct {
	Total int64 `json:"total"` // 数据总数
	Rows  []T   `json:"rows"`  // 数据
}
