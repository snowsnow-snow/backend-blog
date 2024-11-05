package utility

import (
	"backend-blog/config"
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"image"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type Exiftool struct {
}

func (r *Exiftool) ReadExif(path, fileName, extension string) *entity.BlogImage {
	fullPath := config.GlobalConfig.File.Path.System + path + Separator + fileName + Point + extension
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
	var imgInfo entity.BlogImage
	logger.Info.Printf("ExifTool %s\n", config.GlobalConfig.ExifTool.Path+Separator+"exiftool")
	// 调用 exiftool 并读取图像的所有元数据
	cmd := exec.Command(config.GlobalConfig.ExifTool.Path+Separator+"exiftool", "-j", fullPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		logger.Error.Printf("cmd run error %s\n", err)
		return nil
	}
	// 解析 JSON 格式的输出
	var exifData []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &exifData); err != nil {
		logger.Error.Printf("unmarshal error %s\n", err)
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
	rotate, angle, err := setOrientation(&imgInfo, data)
	if err != nil {
		logger.Error.Printf("%s %s\n", fileName, err)
	}
	setSoftware(&imgInfo, data)
	file, err = os.Open(fullPath)
	if err != nil {
		logger.Error.Println("open file error", err)
		return &imgInfo
	}
	// 压缩图片
	compressPaths, err := CompressImage(path, fileName, extension, file)
	if err != nil {
		logger.Error.Printf("%s compress err. %s \n", fileName, err)
		return &imgInfo
	}

	err = RotatePicture(rotate, angle, compressPaths...)
	if err != nil {
		logger.Error.Printf("%s rotate picture err %s \n", fileName, err)
		return &imgInfo
	}

	if extension == "heic" {
		err = JpgToWebp(config.GlobalConfig.File.Path.System+path+Separator+fileName+Point+"jpg",
			config.GlobalConfig.File.Path.System+path+Separator+fileName+Point+"webp",
			imgInfo.ImageWidth,
			imgInfo.ImageHeight)
	} else {
		err = JpgToWebp(fullPath,
			config.GlobalConfig.File.Path.System+path+Separator+fileName+Point+"webp",
			imgInfo.ImageWidth,
			imgInfo.ImageHeight)
	}
	if err != nil {
		logger.Error.Printf("%s jpg to webp err %s \n", fileName, err)
		return &imgInfo
	}
	return &imgInfo
}

func (r *Exiftool) ReadFujiInfo(imgInfo entity.BlogImage, exifData map[string]interface{}) {
}
func handleSpecialParameters(imgInfo *entity.BlogImage, exifData map[string]interface{}) {
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
func setWidthANdHeightDimensionByBounds(imgInfo *entity.BlogImage, file *os.File) {
	img, _, err := image.Decode(file)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	bounds := img.Bounds()
	imgInfo.ImageWidth = bounds.Dx()
	imgInfo.ImageHeight = bounds.Dy()
}
func setOrientation(imgInfo *entity.BlogImage, exifData map[string]interface{}) (bool, int, error) {
	orientation := exifData["Orientation"]
	if orientation == nil {
		return false, 0, nil
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
		fromString := getNumbersFromString(orientationStr)
		return true, int(fromString[0]), nil
	}
	return false, 0, nil
}

func setSoftware(imgInfo *entity.BlogImage, exifData map[string]interface{}) {
	software := exifData["Software"]
	if software == nil {
		return
	}
	valType := reflect.TypeOf(software)
	if valType.Kind() == reflect.String {
		imgInfo.Software = software.(string)
		return
	}
	imgInfo.Software = fmt.Sprintf("%f", software)
	for strings.HasSuffix(imgInfo.Software, "0") {
		imgInfo.Software = strings.TrimSuffix(imgInfo.Software, "0")
	}
	if strings.HasSuffix(imgInfo.Software, ".") {
		imgInfo.Software = strings.TrimSuffix(imgInfo.Software, ".")
	}
	if imgInfo.Make == "Apple" {
		imgInfo.Software = "iOS " + imgInfo.Software
	}
}

func latitudeOrLongitudeStrToCoordinates(imgInfo *entity.BlogImage, exifData map[string]interface{}) error {
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
