package util

import "backend-blog/models"

// ImgUtil 读取图片中的信息
type ImgUtil interface {
	// ReadExif 读取图片 EXIF 信息的通用方法
	ReadExif(path string, fileName string, fileType string) *models.ImgInfo
	// ReadFujiInfo 读取富士相机独有的参数
	ReadFujiInfo(imgInfo models.ImgInfo, exifData map[string]interface{})
}
