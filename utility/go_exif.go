package utility

import (
	"backend-blog/internal/logger"
	"backend-blog/internal/model/entity"
	"errors"
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/tiff"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"runtime/debug"
	"strconv"
)

type GoExif struct{}

func (_ *GoExif) ReadExif(path string, fileName string, fileType string) *entity.BlogImage {
	fullPath := path + Separator + fileName + Point + fileType
	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		logger.Error.Println("Open file error", err)
		return nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			logger.Error.Println("Close file error", err)
		}
	}(file)
	imgInfo := &entity.BlogImage{}
	// 图片宽高信息
	//setPixelXAndYDimensionByBounds(imgInfo, file)
	// 解码 EXIF 数据
	exifInfos, err := exif.Decode(file)
	if err != nil {
		file, err = os.Open(fullPath)
		setPixelXAndYDimensionByBounds(imgInfo, file)
		//err = CompressImage(path, fileName, fileType, 1, file)
		err = errors.New("")
		if err != nil {
			logger.Error.Println("CompressImage file error", err)
		}
		logger.Error.Println("Read exif info error: ", err, string(debug.Stack()))
		return imgInfo
	}
	orientation := ReadOrientation(exifInfos)
	//err = CompressImage(path, fileName, fileType, orientation, file)
	err = errors.New("")
	if err != nil {
		logger.Error.Println("CompressImage file error", err)
	}
	// 图片宽高信息
	//var imgInfo *entity.BlogImage
	//if imgInfo == nil {
	//	logger.Error.Printf("", imgInfo)
	//}
	//imgInfo.Id = uuid.NewString()
	//
	setPixelXAndYDimension(imgInfo, exifInfos, orientation)
	imgInfo.Model = getStringVal(exifInfos, exif.Model)
	imgInfo.Make = getStringVal(exifInfos, exif.Make)
	imgInfo.Software = getStringVal(exifInfos, exif.Software)
	//imgInfo.LensModel = getStringVal(exifInfos, exif.LensModel)
	imgInfo.DateTimeOriginal = getStringVal(exifInfos, exif.DateTimeOriginal)

	// 光圈信息
	setApertureValue(imgInfo, exifInfos)
	// 曝光时间
	setExposureTime(imgInfo, exifInfos)
	// 快门速度
	setShutterSpeedValue(imgInfo, exifInfos)
	// 经纬度信息
	setLatitudeAndLongitude(imgInfo, exifInfos)

	// 快门速度
	setFNumber(imgInfo, exifInfos)
	// 焦距
	setFocalLength(imgInfo, exifInfos)
	setExposureProgram(imgInfo, exifInfos)
	setISO(imgInfo, exifInfos)
	return imgInfo
}
func (_ *GoExif) ReadFujiInfo(imgInfo entity.BlogImage, exifData map[string]interface{}) {
}
func setISO(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	isoProgram, err := exifInfos.Get(exif.ISOSpeedRatings)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	iso, err := isoProgram.Int64(0)
	imgInfo.ISO = iso
}
func ReadOrientation(exifInfos *exif.Exif) int {
	orientation, err := exifInfos.Get(exif.Orientation)
	if err != nil {
		fmt.Println("failed to get orientation, err: ", err)
		return 0
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		fmt.Println("failed to convert type of orientation, err: ", err)
		return 0
	}
	return orientVal
}

func setExposureProgram(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	exposureProgram, err := exifInfos.Get(exif.ExposureProgram)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	dataType := exposureProgram.Type
	info := entity.ExposureProgramMap[uint16(dataType)]
	imgInfo.ExposureProgram = info.ExposureProgram
	imgInfo.ExposureProgramZhCN = info.ExposureProgramZhCN
}

func setFocalLength(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	focalLength, err := exifInfos.Get(exif.FocalLength)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	focalLengthNum, focalLengthDen, err := focalLength.Rat2(0)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imgInfo.FocalLength = formatFloat(float64(focalLengthNum)/float64(focalLengthDen), 2)
}
func setPixelXAndYDimensionByBounds(imgInfo *entity.BlogImage, file *os.File) {
	img, _, err := image.Decode(file)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	bounds := img.Bounds()
	imgInfo.ImageWidth = bounds.Dx()
	imgInfo.ImageHeight = bounds.Dy()
}

func setPixelXAndYDimension(imgInfo *entity.BlogImage, exifInfos *exif.Exif, orientation int) {
	imageWidth, err := exifInfos.Get(exif.PixelXDimension)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	pixelXDimensionVal, err := imageWidth.Int64(0)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imageHeight, err := exifInfos.Get(exif.PixelYDimension)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	pixelYDimensionVal, err := imageHeight.Int64(0)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	//resolutionUnit, _ := exifInfos.Get(exif.ResolutionUnit)
	//resolutionUnitVal, _ := resolutionUnit.Int64(0)
	if orientation == 8 {
		imgInfo.ImageWidth = int(pixelYDimensionVal)
		imgInfo.ImageHeight = int(pixelXDimensionVal)
	} else {
		imgInfo.ImageWidth = int(pixelXDimensionVal)
		imgInfo.ImageHeight = int(pixelYDimensionVal)
	}
}
func setShutterSpeedValue(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	//shutterSpeedValue, err := exifInfos.Get(exif.ShutterSpeedValue)
	//if err != nil {
	//	logger.Warn.Println(err.Error())
	//	return
	//}
	//shutterSpeedValueNum, shutterSpeedValueDen, err := shutterSpeedValue.Rat2(0)
	//if err != nil {
	//	logger.Warn.Println(err.Error())
	//	return
	//}
	//imgInfo.ShutterSpeedValue = float64(shutterSpeedValueNum) / float64(shutterSpeedValueDen)
}

func setExposureTime(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	//exposureTime, err := exifInfos.Get(exif.ExposureTime)
	//if err != nil {
	//	logger.Warn.Println(err.Error())
	//	return
	//}
	//exposureTimeNum, exposureTimeDen, err := exposureTime.Rat2(0)
	//if err != nil {
	//	logger.Warn.Println(err.Error())
	//	return
	//}
	//if sameDigits(exposureTimeNum, exposureTimeDen) {
	//	imgInfo.ExposureTime = denLeftCalculation(exposureTimeNum, exposureTimeDen)
	//} else {
	//	imgInfo.ExposureTime = float64(exposureTimeNum)/10, 0) + "/" + formatFloat(float64(exposureTimeDen)/10
	//}
}

func setApertureValue(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	apertureValue, err := exifInfos.Get(exif.ApertureValue)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	apertureValueNum, apertureValueDen, err := apertureValue.Rat2(0)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imgInfo.ApertureValue = float64(apertureValueNum) / float64(apertureValueDen)
}
func setFNumber(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	fNumber, err := exifInfos.Get(exif.FNumber)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	fNumberNum, fNumberDen, err := fNumber.Rat2(0)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imgInfo.FNumber = float64(fNumberNum) / float64(fNumberDen)
}

func getStringVal(exifInfos *exif.Exif, fieldName exif.FieldName) string {
	tag, err := exifInfos.Get(fieldName)
	if err != nil {
		logger.Warn.Println(err.Error())
		return ""
	}
	val, _ := tag.StringVal()
	return val
}
func setLatitudeAndLongitude(imgInfo *entity.BlogImage, exifInfos *exif.Exif) {
	// 获取经度
	gpsLongitude, err := exifInfos.Get(exif.GPSLongitude)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	// 获取维度
	gpsLatitude, err := exifInfos.Get(exif.GPSLatitude)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	gpsLongitudeRef, err := exifInfos.Get(exif.GPSLongitudeRef)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	gpsLatitudeRef, err := exifInfos.Get(exif.GPSLatitudeRef)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	gpsLongitudeRefVal, err := gpsLongitudeRef.StringVal()
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	gpsLatitudeRefVal, err := gpsLatitudeRef.StringVal()
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	longitudeCoordinate, err := latitudeOrLongitudeToCoordinates(gpsLongitude, gpsLongitudeRefVal)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imgInfo.LongitudeCoordinate = longitudeCoordinate
	latitudeCoordinate, err := latitudeOrLongitudeToCoordinates(gpsLatitude, gpsLatitudeRefVal)
	if err != nil {
		logger.Warn.Println(err.Error())
		return
	}
	imgInfo.LatitudeCoordinate = latitudeCoordinate
}

// latitudeOrLongitudeToCoordinates 经纬度转坐标
func latitudeOrLongitudeToCoordinates(latitudeOrLongitude *tiff.Tag, latitudeRef string) (string, error) {
	var coordinate float64
	var err error
	if degrees, minutes, seconds, gpsLatitudeErr := getDegreesMinutesSeconds(latitudeOrLongitude); gpsLatitudeErr != nil {
		return "0", errors.New("get degrees minutes seconds error")
	} else {
		coordinate, err = dmsToDecimal(degrees, minutes, seconds, latitudeRef)
	}
	if err != nil {
		return "0", errors.New("get coordinate error")
	}
	return formatFloat(coordinate, 20), nil
}

// getDegreesMinutesSeconds 获取度分秒
func getDegreesMinutesSeconds(tag *tiff.Tag) (float64, float64, float64, error) {
	degreesNum, degreesDen, degreesErr := tag.Rat2(0)
	minutesNum, minutesDen, minutesErr := tag.Rat2(1)
	secondsNum, secondsDen, secondsErr := tag.Rat2(2)
	if degreesErr != nil {
		return 0, 0, 0, errors.New("get degreesNum and degreesDen error")
	}
	if minutesErr != nil {
		return 0, 0, 0, errors.New("get minutesNum and minutesDen error")
	}
	if secondsErr != nil {
		return 0, 0, 0, errors.New("secondsNum and secondsDen get err")
	}
	return float64(degreesNum) / float64(degreesDen),
		float64(minutesNum) / float64(minutesDen),
		float64(secondsNum) / float64(secondsDen),
		nil
}

// dmsToDecimal 将度分秒表示的经纬度转换为十进制度数
func dmsToDecimal(degrees float64, minutes float64, seconds float64, direction string) (float64, error) {
	if direction != "N" && direction != "S" && direction != "E" && direction != "W" {
		return 0, errors.New("invalid direction:" + direction)
	}

	sign := 1.0
	if direction == "S" || direction == "W" {
		sign = -1.0
	}

	decimal := sign * (degrees + minutes/60.0 + seconds/3600.0)

	return decimal, nil
}
func formatFloat(num float64, decimal int) string {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}
func sameDigits(num, den int64) bool {
	strNum := strconv.FormatInt(num, 10)
	strDen := strconv.FormatInt(den, 10)
	return len(strNum) == len(strDen)
}
func denLeftCalculation(num, den int64) string {
	strNum := strconv.FormatInt(num, 10)
	floatDen := float64(den)
	for i, n := 0, len(strNum); i < n; i++ {
		floatDen = floatDen / 10
	}
	return strconv.FormatFloat(floatDen, 'f', -1, 64)
}
