package entity

type BlogContent struct {
	BaseInfo
	Title           string `form:"title" json:"title"`                     // 标题
	Text            string `form:"text" json:"text"`                       // 文本
	PublishLocation string `form:"publishLocation" json:"publishLocation"` // 发布位置
	State           int    `form:"state" json:"state"`                     // 1-发布，2-草稿，3-隐藏
	Type            int    `form:"type" json:"type"`                       // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
	Tag             int    `form:"tag" json:"tag"`                         // 标签
	TheCover        string `form:"theCover" json:"theCover"`               // 封面
}

func (BlogContent) TableName() string {
	//实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为articles（结构体+s）
	return "content"
}

//type PublicContentInfo struct {
//	ID              string           `form:"id" json:"id"`                           // ID
//	Title           string           `form:"title" json:"title"`                     // 标题
//	Text            string           `form:"text" json:"text"`                       // 文本
//	PublishLocation string           `form:"publishLocation" json:"publishLocation"` // 发布位置
//	State           int              `form:"state" json:"state"`                     // 1-发布，2-草稿，3-隐藏
//	Type            int              `form:"type" json:"type"`                       // 内容类型：1-文字，2-Markdown，3-图片，4-视频，5-图片和视频
//	Tag             int              `form:"tag" json:"tag"`                         // 标签
//	TheCover        string           `form:"theCover" json:"theCover"`               // 创建时间
//	CreatedTime     common.TimeStamp `json:"createdTime"`
//}
//
//func CreateContent(content BlogContent, c *fiber.Ctx) error {
//	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
//	err := transactionDB.Table(constant.Table.BlogContent).Create(&content)
//	if err != nil {
//		return err.Error
//	}
//	return nil
//}
//
//func DeleteContent(id string, db *gorm.DB) error {
//	err := db.Remove(&BlogContent{BaseInfo: BaseInfo{ID: id}})
//	if err != nil {
//		return err.Error
//	}
//	return nil
//}
//
//func UpdateContent(content BlogContent, c *fiber.Ctx) error {
//	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
//	return transactionDB.Model(&BlogContent{}).
//		Where("id = ?", content.ID).
//		Updates(content).
//		Error
//}
//
//func ListContent(db *gorm.DB) (*[]BlogContent, error) {
//	list := new([]BlogContent)
//	tx := db.Find(&list)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return list, nil
//}
//func PublicListContent(db *gorm.DB) (*[]PublicContentInfo, error) {
//	list := new([]PublicContentInfo)
//	tx := db.Find(&list)
//	if tx.Error != nil {
//		return nil, tx.Error
//	}
//	return list, nil
//}
//func SelectContentById(id string, db *gorm.DB) (*BlogContent, error) {
//	var contentInfo BlogContent
//	err := db.Table(constant.Table.ContentInfo).Where("id = ?", id).Scan(&contentInfo).Error
//	if err != nil {
//		return nil, err
//	}
//	return &contentInfo, nil
//}
//func SelectPublicContentById(id string, db *gorm.DB) (*PublicContentInfo, error) {
//	var contentInfo PublicContentInfo
//	err := db.Table(constant.Table.ContentInfo).Where("id = ?", id).Scan(&contentInfo).Error
//	if err != nil {
//		return nil, err
//	}
//	return &contentInfo, nil
//}
//
//func Count(db *gorm.DB) (int64, error) {
//	var count int64
//	if err := db.Model(BlogContent{}).Count(&count).Error; err != nil {
//		return -1, err
//	}
//	return count, nil
//}
