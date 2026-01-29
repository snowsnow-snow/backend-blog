package media

import (
	"backend-blog/internal/model/entity"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log/slog"
	"os"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"gorm.io/datatypes"
)

func init() {
	// Register mknote readers
	exif.RegisterParsers(mknote.All...)
}

type GoExif struct{}

const (
	Separator = string(os.PathSeparator)
	Point     = "."
)

func (_ *GoExif) ReadExif(fullPath string) *entity.MediaAsset {
	// 打开文件
	file, err := os.Open(fullPath)
	if err != nil {
		slog.Error("Open file error", "error", err)
		return nil
	}
	defer file.Close()

	asset := &entity.MediaAsset{}

	// Decode Config first for basic info
	file.Seek(0, 0)
	cfg, _, err := image.DecodeConfig(file)
	if err == nil {
		asset.Width = cfg.Width
		asset.Height = cfg.Height
	} else {
		slog.Warn("Decode image config failed", "error", err)
	}

	// Reset for EXIF
	file.Seek(0, 0)

	// 解码 EXIF 数据
	exifInfos, err := exif.Decode(file)
	if err != nil {
		slog.Warn("Read exif info error or no exif", "error", err)
		return asset
	}

	// Orientation
	orientation := ReadOrientation(exifInfos)
	if orientation == 6 || orientation == 8 {
		// Swap width/height
		asset.Width, asset.Height = asset.Height, asset.Width
	}

	asset.DeviceMake = getStringVal(exifInfos, exif.Make)
	asset.DeviceModel = getStringVal(exifInfos, exif.Model)

	// Metadata map
	metadata := make(map[string]interface{})

	// Standard Exif
	metadata["software"] = getStringVal(exifInfos, exif.Software)
	metadata["date_time_original"] = getStringVal(exifInfos, exif.DateTimeOriginal)

	// Lens
	// LensModel is not in standard exif package constants mostly, need to check if available or use string lookup
	// goexif doesn't have LensModel constant? It might be tag 0xA434
	if val, err := exifInfos.Get(exif.FieldName("LensModel")); err == nil {
		metadata["lens_model"], _ = val.StringVal()
	}

	// Exposure
	populateExposure(metadata, exifInfos)

	// GPS
	gpsData := getGPS(exifInfos)
	if gpsData != nil {
		metadata["gps"] = gpsData
	}

	// Fujifilm Specific
	if strings.Contains(strings.ToUpper(asset.DeviceMake), "FUJIFILM") {
		ReadFujiInfo(metadata, exifInfos)
	}

	// Marshal metadata
	if metaBytes, err := json.Marshal(metadata); err == nil {
		asset.Metadata = datatypes.JSON(metaBytes)
	}

	return asset
}

func ReadFujiInfo(metadata map[string]interface{}, x *exif.Exif) {
	// Try to get tags. Note: goexif field names from mknote might be specific.
	// We try common names.

	// Film Mode
	if val, err := x.Get(exif.FieldName("FilmMode")); err == nil {
		strVal, _ := val.StringVal()
		metadata["film_mode"] = GetChineseFilmMode(strings.Trim(strVal, "\x00"))
	}

	// Dynamic Range
	if val, err := x.Get(exif.FieldName("DynamicRange")); err == nil {
		strVal, _ := val.StringVal()
		metadata["dynamic_range"] = GetChineseDynamicRange(strings.Trim(strVal, "\x00"))
	}

	// White Balance
	if val, err := x.Get(exif.FieldName("WhiteBalance")); err == nil {
		strVal, _ := val.StringVal()
		metadata["white_balance"] = GetChineseWhiteBalance(strings.Trim(strVal, "\x00"))
	}

	// Sharpness -> Generic description
	if val, err := x.Get(exif.FieldName("Sharpness")); err == nil {
		strVal, _ := val.StringVal()
		metadata["sharpness"] = GetChineseGenericDescriptionMap(strings.Trim(strVal, "\x00"))
	}
}

func populateExposure(m map[string]interface{}, x *exif.Exif) {
	if val, err := x.Get(exif.ISOSpeedRatings); err == nil {
		if iso, err := val.Int64(0); err == nil {
			m["iso"] = iso
		}
	}
	if val, err := x.Get(exif.FNumber); err == nil {
		num, den, _ := val.Rat2(0)
		if den != 0 {
			m["f_number"] = float64(num) / float64(den)
		}
	}
	if val, err := x.Get(exif.ExposureTime); err == nil {
		// Keep as fraction string "1/250"
		num, den, _ := val.Rat2(0)
		if den != 0 {
			m["exposure_time"] = fmt.Sprintf("%d/%d", num, den)
		}
	}
	if val, err := x.Get(exif.FocalLength); err == nil {
		num, den, _ := val.Rat2(0)
		if den != 0 {
			m["focal_length"] = fmt.Sprintf("%.0fmm", float64(num)/float64(den))
		}
	}
}

func getGPS(x *exif.Exif) map[string]float64 {
	lat, long, err := x.LatLong()
	if err != nil {
		return nil
	}
	return map[string]float64{
		"lat": lat,
		"lng": long,
	}
}

func ReadOrientation(x *exif.Exif) int {
	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 0
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		return 0
	}
	return orientVal
}

func getStringVal(x *exif.Exif, fieldName exif.FieldName) string {
	tag, err := x.Get(fieldName)
	if err != nil {
		return ""
	}
	val, _ := tag.StringVal()
	// Clean up null bytes
	return strings.Trim(val, "\x00")
}
