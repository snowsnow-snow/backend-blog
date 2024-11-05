package entity

type BlogMarkdown struct {
	BaseInfo
	FileId string `gorm:"column:file_id;comment:'文件 ID'" json:"fileId"`
	Night  string `gorm:"column:night;comment:'是否为黑暗模式,0:否,1:是'" json:"night"` // 是否为黑暗模式
}

//	func CreateMarkdownInfo(ii BlogMarkdown, c *fiber.Ctx) error {
//		transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
//		err := transactionDB.Table(constant.Table.MarkdownInfo).Create(&ii)
//		if err != nil {
//			return err.Error
//		}
//		return nil
//	}
func (BlogMarkdown) TableName() string {
	return "markdown"
}
