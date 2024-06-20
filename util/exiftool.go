package util

import (
	"backend-blog/logger"
	"backend-blog/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"image"
	"os"
	"os/exec"
	"strings"
)

type Exiftool struct {
}

func (r *Exiftool) ReadExif(path string, fileName string, fileType string) *models.ImgInfo {
	fullPath := path + Separator + fileName + Point + fileType
	// 打开文件
	file, err := os.Open(fullPath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error.Println("read exif close file error", err)
		}
	}(file)
	if err != nil {
		logger.Error.Println("open file error", err)
		return nil
	}
	var imgInfo models.ImgInfo

	// 调用 exiftool 并读取图像的所有元数据
	cmd := exec.Command("image-ExifTool-12.86/exiftool", "-j", fullPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil
	}
	// 解析 JSON 格式的输出
	var exifData []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &exifData); err != nil {
		logger.Error.Println("unmarshal error", err)
	}
	if len(exifData) == 0 {
		return nil
	}
	data := exifData[0]
	err = mapstructure.Decode(data, &imgInfo)
	handleSpecialParameters(&imgInfo, data)
	if err != nil {
		logger.Error.Printf("%s %s\n", fileName, err)
	}
	err = latitudeOrLongitudeStrToCoordinates(&imgInfo, data)
	if err != nil {
		logger.Error.Printf("%s %s\n", fileName, err)
	}
	err = redAndBlue(&imgInfo, data)
	if err != nil {
		logger.Error.Printf("%s %s\n", fileName, err)
	}
	rotate, err := setOrientation(fullPath, &imgInfo, data)
	if err != nil {
		logger.Error.Printf("%s %s\n", fileName, err)
	}
	file, err = os.Open(fullPath)
	if err != nil {
		logger.Error.Println("open file error", err)
		return &imgInfo
	}
	// 压缩图片
	compressPaths, err := Compress(path, fileName, fileType, file)
	if err != nil {
		logger.Error.Printf("%s compress err.%s \n", fileName, err)
		return &imgInfo
	}
	err = RotatePicture(rotate, compressPaths...)
	if err != nil {
		logger.Error.Printf("%s rotate picture err.%s \n", fileName, err)
		return &imgInfo
	}
	return &imgInfo
}

func (r *Exiftool) ReadFujiInfo(imgInfo models.ImgInfo, exifData map[string]interface{}) {
}
func handleSpecialParameters(imgInfo *models.ImgInfo, exifData map[string]interface{}) {
	if imgInfo.ExposureTime == "" {
		if temp := exifData["ExposureTime"]; temp != nil {
			imgInfo.ExposureTime = fmt.Sprint(temp)
		}
	}
	if imgInfo.ShutterSpeedValue == "" {
		if temp := exifData["ShutterSpeedValue"]; temp != nil {
			imgInfo.ShutterSpeedValue = fmt.Sprint(temp)
		}
	}
	if imgInfo.FocalLength != "" {
		imgInfo.FocalLength = strings.ReplaceAll(imgInfo.FocalLength, " ", "")
		imgInfo.FocalLength = strings.ReplaceAll(imgInfo.FocalLength, ".0", "")
	}
	if imgInfo.FocalLengthIn35mmFormat != "" {
		imgInfo.FocalLengthIn35mmFormat = strings.ReplaceAll(imgInfo.FocalLengthIn35mmFormat, " ", "")
		imgInfo.FocalLengthIn35mmFormat = strings.ReplaceAll(imgInfo.FocalLengthIn35mmFormat, ".0", "")
	}
}
func setWidthANdHeightDimensionByBounds(imgInfo *models.ImgInfo, file *os.File) {
	img, _, err := image.Decode(file)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	bounds := img.Bounds()
	imgInfo.ImageWidth = bounds.Dx()
	imgInfo.ImageHeight = bounds.Dy()
}
func setOrientation(fullPath string, imgInfo *models.ImgInfo, exifData map[string]interface{}) (bool, error) {
	orientation := exifData["Orientation"]
	if orientation == nil {
		return false, nil
	}
	orientationStr, _ := orientation.(string)
	if strings.Contains(orientationStr, "Rotate") {
		//err := RotatePicture90(fullPath)
		//if err != nil {
		//	return true, err
		//}
		oldWidth := imgInfo.ImageWidth
		imgInfo.ImageWidth = imgInfo.ImageHeight
		imgInfo.ImageHeight = oldWidth
		return true, nil
	}
	return false, nil
}
func redAndBlue(imgInfo *models.ImgInfo, exifData map[string]interface{}) error {
	//whiteBalanceFineTune := fmt.Sprint(exifData["WhiteBalanceFineTune"])
	//if whiteBalanceFineTune == "" {
	//	return errors.New("whiteBalanceFineTune or gpsLatitude is empty")
	//}
	//split := strings.Split(whiteBalanceFineTune, ",")
	//red, blue := split[0], split[1]
	//if red != "" {
	//	imgInfo.Red = red[len(red)-2:]
	//}
	//if blue != "" {
	//	imgInfo.Blue = blue[len(blue)-2:]
	//}
	return nil
}
func latitudeOrLongitudeStrToCoordinates(imgInfo *models.ImgInfo, exifData map[string]interface{}) error {
	gpsLongitude := fmt.Sprint(exifData["GPSLongitude"])
	gpsLatitude := fmt.Sprint(exifData["GPSLatitude"])
	if (gpsLongitude == "" || gpsLatitude == "") || (gpsLongitude == "<nil>" || gpsLatitude == "<nil>") {
		return errors.New("gpsLongitude or gpsLatitude is empty")
	}
	gpsLongitudeArray := getNumbersFromString(gpsLongitude)
	gpsLatitudeArray := getNumbersFromString(gpsLatitude)
	finalLongitude, err := dmsToDecimal(gpsLongitudeArray[0], gpsLongitudeArray[1], gpsLongitudeArray[2], string(gpsLongitude[len(gpsLongitude)-1]))
	if err != nil {
		return errors.New("get final longitude error")
	}
	finalLatitude, err := dmsToDecimal(gpsLatitudeArray[0], gpsLatitudeArray[1], gpsLatitudeArray[2], string(gpsLatitude[len(gpsLatitude)-1]))
	if err != nil {
		return errors.New("get final longitude error")
	}

	imgInfo.LongitudeCoordinate = formatFloat(finalLongitude, 5)
	imgInfo.LatitudeCoordinate = formatFloat(finalLatitude, 5)
	return nil
}
