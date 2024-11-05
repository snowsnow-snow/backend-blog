package constant

var (
	Table = tableNames{
		"user",
		"file",
		"content",
		"image",
		"video",
		"markdown",
	}
	Local = localNames{
		"TransactionDB", // 数据库事物
	}
	YesOrNo = yesOrNo{
		1,
		0,
	}
	// ImageCompressionRatioMap 压缩比例
	ImageCompressionRatioMap = map[string]string{
		"50":  "50",
		"70":  "70",
		"100": "100",
	}
)

type tableNames struct {
	User         string
	File         string
	BlogContent  string
	BlogImage    string
	BlogVideo    string
	BlogMarkdown string
}
type localNames struct {
	TransactionDB string
}

type yesOrNo struct {
	Yes int
	No  int
}

// ImageExtensions 支持的图像扩展名
var ImageExtensions = []string{
	"jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp",
}

// VideoExtensions 支持的视频扩展名
var VideoExtensions = []string{
	"mp4", "avi", "mov", "wmv", "mkv", "flv", "webm",
}
