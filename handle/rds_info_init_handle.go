package handle

import (
	"backend-blog/common"
	"backend-blog/models"
	"backend-blog/util"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Img struct{}
type Video struct{}
type LivePhotos struct{}
type Markdown struct{}

type RdsInfoInitHandle interface {
	InitFileInfo(saveFile *models.FileInfo, c *fiber.Ctx) (*models.SaveFileInfo, error)
}

func (r Img) InitFileInfo(imgInfo *models.FileInfo, c *fiber.Ctx) (*models.SaveFileInfo, error) {
	var exiftool util.Exiftool
	img := exiftool.ReadExif(imgInfo.FilePath, imgInfo.FileName, imgInfo.RawFileName, imgInfo.Type)
	if img == nil {
		img = &models.ImgInfo{}
	}
	common.CreateInit(c, &img.BaseInfo)
	return &models.SaveFileInfo{
		FileInfo: *imgInfo,
		ImgInfo:  *img,
	}, nil
}

func (r Video) InitFileInfo(videoInfo *models.FileInfo, c *fiber.Ctx) (*models.SaveFileInfo, error) {
	video, err := util.ReadInfo(videoInfo.FilePath +
		util.Separator +
		videoInfo.FileName +
		util.Point +
		videoInfo.Type)
	if err != nil {
		return nil, err
	}
	common.CreateInit(c, &video.BaseInfo)
	return &models.SaveFileInfo{
		FileInfo:  *videoInfo,
		VideoInfo: *video,
	}, nil
}

func (r LivePhotos) InitFileInfo(saveFile *models.FileInfo, c *fiber.Ctx) (*models.SaveFileInfo, error) {
	if util.IsInArrayNoCaseSensitive(util.ImgType, saveFile.Type) {
		var img = &Img{}
		return img.InitFileInfo(saveFile, c)
	}
	if util.IsInArrayNoCaseSensitive(util.VideoType, saveFile.Type) {
		var video = &Video{}
		return video.InitFileInfo(saveFile, c)
	}
	return nil, nil
}
func (r Markdown) InitFileInfo(saveFile *models.FileInfo, c *fiber.Ctx) (*models.SaveFileInfo, error) {
	var file models.SaveFileInfo
	file.FileInfo = *saveFile
	markdownInfo := models.MarkdownInfo{}
	if strings.Contains(file.RawFileName, "_night") {
		markdownInfo.Night = "1"
	} else {
		markdownInfo.Night = "0"
	}
	common.CreateInit(c, &markdownInfo.BaseInfo)
	return &models.SaveFileInfo{
		FileInfo:     *saveFile,
		MarkdownInfo: markdownInfo,
	}, nil
}
