package entity

type Category struct {
	BaseInfo
	Name        string `Gorm:"uniqueIndex;size:100;not null" json:"name"`
	Description string `Gorm:"size:255" json:"description"`
}

func (Category) TableName() string {
	return "category"
}
