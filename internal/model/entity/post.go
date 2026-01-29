package entity

type PostType string

const (
	PostTypeArticle PostType = "article" // 混合/普通文章
	PostTypeGallery PostType = "gallery" // 纯图片模式
	PostTypeVideo   PostType = "video"   // 纯视频模式
)

type Post struct {
	BaseInfo
	PostType PostType `Gorm:"default:article;index" json:"postType"`
	Title    string   `Gorm:"size:255" json:"title"`
	Summary  string   `Gorm:"size:500" json:"summary"`  // 一句话描述
	Content  string   `Gorm:"type:text" json:"content"` // Markdown 内容

	// 关联关系
	CoverImageID *int64       `json:"coverImageId,string"` // 封面图ID，可选
	CategoryID   *int64       `Gorm:"index" json:"categoryId,string"`
	Category     *Category    `Gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	MediaAssets  []MediaAsset `Gorm:"foreignKey:PostID;constraint:OnDelete:SET NULL;" json:"mediaAssets"`

	Status string `Gorm:"default:published;index" json:"status"` // published, draft
}

func (Post) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为articles（结构体+s）
	return "post"
}
