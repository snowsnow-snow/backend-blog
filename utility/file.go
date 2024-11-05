package utility

import (
	"backend-blog/config"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/bo"
	"backend-blog/internal/model/entity"
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
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

var (
	VideoType             = []string{"MOV", "MP4", "WEBM"}
	ImgType               = []string{"PNG", "JPG", "JPEG", "HEIC"}
	MarkDownType          = []string{"HTML", "MD"}
	CompressionRatio      = []int{50, 70, 30}
	VideoCompressionRatio = map[string]int{
		"50": 30,
		"70": 23}
	Separator = string(filepath.Separator)
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

func ParamToFileBo(c *fiber.Ctx) ([]*bo.FileBo, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	contentId := form.Value["contentId"][0]
	sort, _ := strconv.Atoi(form.Value["sort"][0])
	coverParam := form.Value["cover"]
	var cover int
	if coverParam != nil {
		cover, _ = strconv.Atoi(form.Value["cover"][0])
	}
	// 获取所有上传的文件
	files := form.File["files"]
	fileBos := make([]*bo.FileBo, len(files))
	for index, file := range files {
		fileDo, err := initFileDo(file, contentId)
		fileDo.File = file
		if err != nil {
			return nil, err
		}
		fileDo.Sort = sort
		fileDo.Cover = cover
		fileBos[index] = fileDo
	}
	return fileBos, nil
}

// SaveFiles 将上传的文件保存在本地，保存在服务器的文件名字为 time.Now()
// 文件保存路径在 config 文件中配置 config.GlobalConfig.File.Path.Public
func SaveFiles(c *fiber.Ctx, fileBos []*bo.FileBo) ([]*bo.ResourceDescBo, error) {
	resourceDescDos := make([]*bo.ResourceDescBo, len(fileBos))
	for i := range fileBos {
		file, err := SaveFile(c, fileBos[i])
		if err != nil {
			return nil, err
		}
		resourceDescDos[i] = file
	}
	return resourceDescDos, nil

}

// SaveFile 将文件保存
func SaveFile(c *fiber.Ctx, fileBo *bo.FileBo) (*bo.ResourceDescBo, error) {
	// 创建文件夹
	err := CreateFileFolder(config.GlobalConfig.File.Path.System + fileBo.FilePath)
	if err != nil {
		return nil, err
	}
	filePathAndName := config.GlobalConfig.File.Path.System + fileBo.FilePath + fileBo.FileName + Point + fileBo.Extension
	// 将文件保存在本地
	err = c.SaveFile(fileBo.File, filePathAndName)
	if err != nil {
		return nil, err
	}
	// 如果是 heic 类型的文件，将 heic 转为 jpg 文件格式
	if fileBo.Extension == "heic" {
		err := ConvertHeicToJpg(filePathAndName,
			config.GlobalConfig.File.Path.System+fileBo.FilePath+fileBo.FileName+Point+"jpg")
		if err != nil {
			return nil, err
		}
		fileBo.Type = "jpg"
	}
	return &bo.ResourceDescBo{
		File: entity.File{
			RawFileName: fileBo.RawFileName,
			FilePath:    fileBo.FilePath,
			FileName:    fileBo.FileName,
			Type:        fileBo.Type,
			Extension:   fileBo.Extension,
			Size:        fileBo.Size,
			ContentId:   fileBo.ContentId,
			Sort:        fileBo.Sort,
			Cover:       fileBo.Cover,
		},
	}, err
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
		logger.Error.Printf("filename: %s rotate picture read image file error %s\n", path, err)
		return err
	}
	return nil
}
func CompressImage(path, fileName, extension string, currFile *os.File) ([]string, error) {
	path = config.GlobalConfig.File.Path.System + path
	if extension == "heic" {
		extension = "jpg"
	}
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
			logger.Error.Printf("Rotate picture close file error %s\n", err)
		}
	}(file)
	for i := range CompressionRatio {
		var opt jpeg.Options
		opt.Quality = 1
		// 设置压缩质量
		options := &jpeg.Options{Quality: CompressionRatio[i]}
		compressName := fileName + Delimiter + strconv.Itoa(CompressionRatio[i]) + Point + extension
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
func CompressVideo(path, fileName, fileType string) error {
	path = config.GlobalConfig.File.Path.System + path
	for key, value := range VideoCompressionRatio {
		// 构造 FFmpeg 命令
		cmd := exec.Command("/Users/snowsnowsnow/MyBlog/ffmpeg/ffmpeg",
			"-i", path+fileName+Point+fileType, "-vcodec", "libx264", "-crf",
			fmt.Sprintf("%d", value), path+fileName+Delimiter+key+Point+fileType)

		// 执行命令
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
func JpgToWebp(jpgPath, webpPath string, width, height int) error {
	// 缩小比例
	var scaleFactor float64
	if width == 1290 || height == 2796 {
		scaleFactor = 1
	} else {
		scaleFactor = 0.25
	}
	// 压缩质量
	quality := 100
	// 使用 FFmpeg 将 JPEG 转换为 WebP
	cmd := exec.Command(
		"/Users/snowsnowsnow/MyBlog/ffmpeg/ffmpeg",
		"-i", jpgPath,
		"-vf", fmt.Sprintf("scale=iw*%f:ih*%f", scaleFactor, scaleFactor),
		"-q:v", fmt.Sprintf("%d", quality),
		webpPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		logger.Error.Printf("转换失败: %v", err)
		return err
	}
	return nil
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFiles(paths ...string) error {
	for _, path := range paths {
		err := DeleteFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}
func DeleteFolder(folderPath string) error {
	err := os.RemoveAll(folderPath)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFolders(folderPaths ...string) error {
	for _, folderPath := range folderPaths {
		err := DeleteFolder(folderPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func initFileDo(file *multipart.FileHeader, contentId string) (*bo.FileBo, error) {
	filenameArray := strings.Split(file.Filename, Point)
	extension := strings.ToLower(filenameArray[len(filenameArray)-1])
	timeStr := time.Now().Format("2006-01-02")
	filePath := config.GlobalConfig.File.Path.Resource + Separator
	if IsInArray(VideoType, strings.ToUpper(extension)) {
		filePath += "video"
	} else if IsInArray(ImgType, strings.ToUpper(extension)) {
		filePath += "img"
	} else if IsInArray(MarkDownType, strings.ToUpper(extension)) {
		filePath += "markdown"
	}
	timeUnix := fmt.Sprintf("%v", time.Now().UnixNano()/int64(time.Millisecond))
	filePath += Separator + timeStr + Separator + contentId + Separator + timeUnix + Separator
	return &bo.FileBo{
		RawFileName: file.Filename,
		FilePath:    filePath,
		FileName:    timeUnix,
		Type:        file.Header.Get("Content-Type"),
		Extension:   extension,
		Size:        file.Size,
		ContentId:   contentId,
	}, nil
}
