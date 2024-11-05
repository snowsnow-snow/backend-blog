package handle

import (
	"backend-blog/internal/common"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/bo"
	entity "backend-blog/internal/model/entity"
	"backend-blog/utility"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Img struct{}
type Video struct{}
type LivePhotos struct{}
type Markdown struct{}

type RdsInfoInitHandle interface {
	InitFileInfo(saveFile *entity.File, c *fiber.Ctx)
}

func (r Img) InitFileInfo(resourceDescDo *bo.ResourceDescBo, c *fiber.Ctx) {
	var exiftool utility.Exiftool
	img := exiftool.ReadExif(resourceDescDo.File.FilePath,
		resourceDescDo.File.FileName,
		resourceDescDo.File.Extension)
	if img == nil {
		img = &entity.BlogImage{}
	}
	common.CreateInit(c, &img.BaseInfo)
	resourceDescDo.BlogImage = *img

}

func (r Video) InitFileInfo(resourceDescDo *bo.ResourceDescBo, c *fiber.Ctx) {
	video, err := utility.ReadInfo(resourceDescDo.File.FilePath +
		resourceDescDo.File.FileName +
		utility.Point +
		resourceDescDo.File.Extension)
	if err != nil {
		logger.Error.Println(err)
	}
	// 压缩视频
	err = utility.CompressVideo(resourceDescDo.File.FilePath, resourceDescDo.File.FileName, resourceDescDo.File.Extension)
	if err != nil {
		logger.Error.Printf("compress video err%s\n", err)
	}
	common.CreateInit(c, &video.BaseInfo)
	resourceDescDo.BlogVideo = *video
}

func (r LivePhotos) InitFileInfo(resourceDescDo *bo.ResourceDescBo, c *fiber.Ctx) {
	if utility.IsInArrayNoCaseSensitive(utility.ImgType, resourceDescDo.File.Extension) {
		var img = &Img{}
		img.InitFileInfo(resourceDescDo, c)
	}
	if utility.IsInArrayNoCaseSensitive(utility.VideoType, resourceDescDo.File.Extension) {
		var video = &Video{}
		video.InitFileInfo(resourceDescDo, c)
	}
}

func (r Markdown) InitFileInfo(rd *bo.ResourceDescBo, c *fiber.Ctx) {
	markdownInfo := entity.BlogMarkdown{}
	if strings.Contains(rd.File.RawFileName, "_night") {
		markdownInfo.Night = "1"
	} else {
		markdownInfo.Night = "0"
	}
	common.CreateInit(c, &markdownInfo.BaseInfo)
	rd.BlogMarkdown = markdownInfo

}
