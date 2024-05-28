package controller

import (
	"backend-blog/logger"
	"backend-blog/services"
	"backend-blog/util"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"

	"backend-blog/config"
	"backend-blog/result"
)

// 压缩比例
var compressionRatioMap = map[string]string{
	"50":  "50",
	"70":  "70",
	"100": "100",
}

type fileController struct {
}

var FileController = new(fileController)

func (r fileController) ViewImage(c *fiber.Ctx) error {
	compressionRatio := c.Params("compressionRatio")
	if compressionRatioMap[compressionRatio] == "" {
		return result.FailWithMsg(c, result.WrongParameter.Error())
	}
	path, _ := services.FileService.ByIdCompressionRatioGetImgPath(c.Params("imageId"), compressionRatio)
	return c.SendFile(path)
}
func (r fileController) ViewVideo(c *fiber.Ctx) error {
	_, proportion := c.Params("isRaw"), c.Params("proportion")
	if compressionRatioMap[proportion] == "" {
		return result.FailWithMsg(c, result.WrongParameter.Error())
	}
	path, _ := services.FileService.ByIdGetVideoPath(c.Params("videoId"))
	return c.SendFile(path)
}
func (r fileController) UploadImg(c *fiber.Ctx) error {
	services.FileService.SaveFile(c)
	//filePath, newFileName, uploadFileName, err := saveImg(c)
	//if err != nil {
	//	logger.Error.Println("save file error:", err)
	//	return result.ErrorWithMsg(c, "save file error")
	//}
	//
	//info, err := readImgInfo(filePath, newFileName, uploadFileName)
	//if err != nil {
	//	logger.Error.Println("read file info error:", err)
	//	return result.ErrorWithMsg(c, "read file info error")
	//}
	//transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	//transactionDB.Begin()
	//createErr := transactionDB.Table(constant.Table.ImgInfo).Create(info)
	//if createErr.Error != nil {
	//	logger.Error.Println("create img info error:", createErr.Error.Error())
	//	transactionDB.Rollback()
	//	return result.Error(c)
	//}
	//transactionDB.Commit()
	return result.Success(c)
}

// saveImg 将上传的文件保存在本地，保存在服务器的文件名字为 time.Now()
// 文件保存路径在 config 文件中配置 config.GlobalConfig.File.Path.Public
func saveImg(c *fiber.Ctx) (string, string, string, error) {
	file, err := c.FormFile("files")
	if err != nil {
		logger.Error.Println("get file error %s\n", err)
		return "", "", "", err
	}
	timeStr := time.Now().Format("2006-01-02")
	filePath := config.GlobalConfig.File.Path.Public + util.Separator + timeStr
	err = util.CreateFileFolder(filePath)
	if err != nil {
		logger.Error.Println("create file older error %s\n", err)
		return "", "", "", err
	}
	timeUnix := time.Now().Unix()
	splitRawFileName := strings.Split(file.Filename, ".")
	newFileName := strconv.FormatInt(timeUnix, 12) + "." + splitRawFileName[len(splitRawFileName)-1]
	filePathAndName := filePath + util.Separator + newFileName
	err = c.SaveFile(file, filePathAndName)
	if err != nil {
		return "", "", "", err
	}
	return filePath, newFileName, file.Filename, nil
}

// readImgInfo 读取图片的 EXIF 信息
//func readImgInfo(filePath, newFileName, uploadFileName string) (*models.ImgInfo, error) {
//	logger.Info.Println("img:", uploadFileName, "read exif start, file path:", filePath, "file upload name:", uploadFileName)
//	exif, err := util.ReadExif(filePath + separator + newFileName)
//	if err != nil {
//		return nil, err
//	}
//	if exif == nil {
//		return nil, errors.New("exif is empty")
//	}
//	logger.Info.Println("img:", uploadFileName, "read exif end, exif:", exif)
//	return exif, nil
//}

func SetIsFIle(c *fiber.Ctx) error {
	c.Locals("isFile", true)
	return c.Next()
}
