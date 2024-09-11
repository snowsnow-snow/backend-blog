package constant

var (
	Table = TableNames{
		"user",
		"img_info",
		"video_info",
		"content_info",
		"resource_desc",
		"markdown_info",
	}
	Local = LocalNames{
		"TransactionDB", // 数据库事物
	}
)

type TableNames struct {
	User         string
	ImgInfo      string
	VideoInfo    string
	ContentInfo  string
	ResourceDesc string
	MarkdownInfo string
}
type LocalNames struct {
	TransactionDB string
}
