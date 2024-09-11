package util

import (
	"backend-blog/config"
	"backend-blog/logger"
	"backend-blog/models"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/gofiber/fiber/v2"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var (
	VideoType        = []string{"MOV", "MP4", "WEBM"}
	ImgType          = []string{"PNG", "JPG", "JPEG", "HEIC"}
	MarkDownType     = []string{"HTML", "MD"}
	CompressionRatio = []int{50, 70, 1}
	Separator        = string(filepath.Separator)
)

func FolderExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	return os.IsExist(err)
}
func CreateFolder(folderPath string) error {
	var err error
	if FolderNotExist(folderPath) {
		err = os.Mkdir(folderPath, 0755)
	}
	return err
}
func FolderNotExist(folderPath string) bool {
	_, err := os.Stat(folderPath)
	return os.IsNotExist(err)
}

func CreateFileFolder(folderPath string) error {
	if FolderNotExist(folderPath) {
		err := os.MkdirAll(folderPath, os.ModePerm)
		if err != nil {
			logger.Error.Println(err, string(debug.Stack()))
			return err
		}
	}
	return nil
}

// SaveFiles 将上传的文件保存在本地，保存在服务器的文件名字为 time.Now()
// 文件保存路径在 config 文件中配置 config.GlobalConfig.File.Path.Public
func SaveFiles(c *fiber.Ctx) ([]*models.FileInfo, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	// 获取所有上传的文件
	files := form.File["files"]
	var fileInfos []*models.FileInfo
	for i := range files {
		file, err := SaveFile(c, files[i])
		fileInfos = append(fileInfos, file)
		if err != nil {
			return nil, err
		}
		//if IsInArray(ImgType, strings.ToUpper(file.Type)) {
		//	err = Compress(*file)
		//}
		if err != nil {
			return nil, err
		}
	}
	return fileInfos, nil

}
func SaveFile(c *fiber.Ctx, file *multipart.FileHeader) (*models.FileInfo, error) {
	fileInfo, err := getFileInfoByFileHeader(file)
	if err != nil {
		return fileInfo, err
	}
	filePathAndName := fileInfo.FilePath + Separator + fileInfo.FileName + Point + fileInfo.Type
	if strings.ToLower(fileInfo.Type) == "heic" {
		err = c.SaveFile(file, filePathAndName)
		if err != nil {
			return fileInfo, err
		}
		//err := os.Rename(filePathAndName, fileInfo.FilePath+Separator+fileInfo.FileName+Point+"jpg")
		err := ConvertHeicToJpg(filePathAndName, fileInfo.FilePath+Separator+fileInfo.FileName+Point+"jpg")
		if err != nil {
			return nil, err
		}
		fileInfo.Type = "jpg"
		//fileInfo.RawFileName = fileInfo.FileName + Point + "jpg"
	} else {
		strings.ToLower(fileInfo.Type)
		err = c.SaveFile(file, filePathAndName)
	}
	return fileInfo, err
}
func RotatePicture(rotate bool, angle int, paths ...string) error {
	if !rotate {
		return nil
	}
	for _, path := range paths {
		err := RotatePicture90(path, angle)
		if err != nil {
			return err
		}
	}
	return nil
}

func RotatePicture90(path string, angle int) error {
	// 打开图像文件
	img, err := imaging.Open(path)
	// 旋转图像（顺时针旋转）
	if angle == 90 {
		img = imaging.Rotate270(img)
	} else if angle == 180 {
		img = imaging.Rotate180(img)
	} else if angle == 270 {
		img = imaging.Rotate90(img)
	}
	err = imaging.Save(img, path)
	if err != nil {
		logger.Error.Println("rotate picture read image file error", err)
		return err
	}
	return nil
}
func Compress(path, fileName, fileType string, currFile *os.File) ([]string, error) {
	_, err := currFile.Seek(0, 0)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(currFile)
	if err != nil {
		return nil, err
	}
	compressPaths := make([]string, len(CompressionRatio))
	var file *os.File
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error.Println("Rotate picture close file error", err)
		}
	}(file)
	for i := range CompressionRatio {
		var opt jpeg.Options
		opt.Quality = 1
		// 设置压缩质量
		options := &jpeg.Options{Quality: CompressionRatio[i]}
		compressName := fileName + Delimiter + strconv.Itoa(CompressionRatio[i]) + Point + fileType
		outputFile, err := os.Create(path + Separator + compressName)
		if err != nil {
			return nil, err
		}
		err = jpeg.Encode(outputFile, img, options)
		if err != nil {
			return nil, err
		}
		err = outputFile.Close()
		if err != nil {
			return nil, err
		}
		compressPaths[i] = path + Separator + compressName
	}
	return compressPaths, nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func getFileInfoByFileHeader(file *multipart.FileHeader) (*models.FileInfo, error) {
	filenameArray := strings.Split(file.Filename, Point)
	fileType := strings.ToLower(filenameArray[len(filenameArray)-1])
	timeStr := time.Now().Format("2006-01-02")
	var filePath string
	if IsInArray(VideoType, strings.ToUpper(fileType)) {
		filePath = config.GlobalConfig.File.Path.Public + Separator + "video" + Separator + timeStr
	} else if IsInArray(ImgType, strings.ToUpper(fileType)) {
		filePath = config.GlobalConfig.File.Path.Public + Separator + "img" + Separator + timeStr
	} else if IsInArray(MarkDownType, strings.ToUpper(fileType)) {
		filePath = config.GlobalConfig.File.Path.Public + Separator + "markdown" + Separator + timeStr
	}

	err := CreateFileFolder(filePath)
	if err != nil {
		logger.Error.Println("create file older error", err)
		return &models.FileInfo{
			RawFileName: file.Filename,
			FilePath:    filePath,
		}, err
	}
	timeUnix := time.Now().UnixNano() / int64(time.Millisecond)
	splitRawFileName := strings.Split(file.Filename, Point)
	//newFileName := fmt.Sprintf("%v", timeUnix) + "." + splitRawFileName[len(splitRawFileName)-1]
	return &models.FileInfo{
		RawFileName: file.Filename,
		FilePath:    filePath,
		FileName:    fmt.Sprintf("%v", timeUnix),
		Type:        splitRawFileName[len(splitRawFileName)-1],
		Size:        file.Size,
	}, nil
}
