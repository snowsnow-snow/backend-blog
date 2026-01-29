package utility

import (
	"backend-blog/config"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ExifData represents the extracted metadata from an image.
type ExifData struct {
	Make                 string             `json:"make"`
	Model                string             `json:"model"`
	Width                int                `json:"width"`
	Height               int                `json:"height"`
	ISO                  int                `json:"iso"`
	FNumber              float64            `json:"f_number"`
	ExposureTime         string             `json:"exposure_time"`
	FocalLength          string             `json:"focal_length"`
	LensModel            string             `json:"lens_model"`
	Software             string             `json:"software"`
	DateTime             string             `json:"date_time"`
	ExposureBias         string             `json:"exposure_bias,omitempty"`
	FilmSimulation       string             `json:"film_simulation,omitempty"`
	DynamicRange         string             `json:"dynamic_range,omitempty"`
	WhiteBalance         string             `json:"white_balance,omitempty"`
	WhiteBalanceFineTune string             `json:"white_balance_fine_tune,omitempty"`
	Sharpness            string             `json:"sharpness,omitempty"`
	NoiseReduction       string             `json:"noise_reduction,omitempty"`
	ShadowTone           string             `json:"shadow_tone,omitempty"`
	Saturation           string             `json:"saturation,omitempty"`
	ColorChromeFXBlue    string             `json:"color_chrome_fx_blue,omitempty"`
	ColorChromeEffect    string             `json:"color_chrome_effect,omitempty"`
	GrainEffectRoughness string             `json:"grain_effect_roughness,omitempty"`
	HighlightTone        string             `json:"highlight_tone,omitempty"`
	LivePhotosId         string             `json:"live_photos_id,omitempty"`
	GPS                  map[string]float64 `json:"gps,omitempty"`
	Raw                  []byte             `json:"-"` // Raw JSON from exiftool
}

type Exiftool struct {
}

// ReadExif calls the exiftool CLI to read metadata from the given image path and returns an ExifData struct.
func (r *Exiftool) ReadExif(fullPath string) (*ExifData, error) {
	// 获取 exiftool 路径，如果配置中没有则默认使用系统路径中的 exiftool
	exifToolPath := config.GlobalConfig.ExifTool.Path
	var binPath string
	if exifToolPath == "" {
		binPath = "exiftool"
	} else {
		binPath = filepathJoin(exifToolPath, "exiftool")
	}

	// 调用 exiftool 并读取图像的所有元数据为 JSON
	cmd := exec.Command(binPath, "-j", fullPath)

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("exiftool execution failed: %w", err)
	}

	// 解析 JSON 格式的输出
	var exifData []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &exifData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal exiftool output: %w", err)
	}

	if len(exifData) == 0 {
		return nil, fmt.Errorf("no metadata found in image")
	}

	data := exifData[0]
	res := &ExifData{
		Raw: out.Bytes(),
	}

	// 映射基础字段
	res.Make = r.toString(data["Make"])
	res.Model = r.toString(data["Model"])
	res.Width = r.toInt(data["ImageWidth"])
	if res.Width == 0 {
		res.Width = r.toInt(data["ExifImageWidth"])
	}
	res.Height = r.toInt(data["ImageHeight"])
	if res.Height == 0 {
		res.Height = r.toInt(data["ExifImageHeight"])
	}

	res.ISO = r.toInt(data["ISO"])
	res.FNumber = r.toFloat64(data["FNumber"])
	res.ExposureTime = r.toString(data["ExposureTime"])
	res.ExposureBias = r.toString(data["ExposureCompensation"])
	res.FocalLength = r.toString(data["FocalLength"])
	res.LensModel = r.toString(data["LensModel"])
	res.Software = r.toString(data["Software"])
	res.DateTime = r.toString(data["DateTimeOriginal"])

	// 富士胶片模拟及详细参数处理
	if strings.Contains(strings.ToUpper(res.Make), "FUJIFILM") {
		// 富士通常将胶片模拟存储在 FilmMode 或 FilmSimulation 中
		if val, ok := data["FilmMode"]; ok {
			res.FilmSimulation = r.toString(val)
		} else if val, ok := data["FilmSimulation"]; ok {
			res.FilmSimulation = r.toString(val)
		}
		res.DynamicRange = r.toString(data["DynamicRange"])
		res.WhiteBalance = r.toString(data["WhiteBalance"])
		res.WhiteBalanceFineTune = r.toString(data["WhiteBalanceFineTune"])
		res.Sharpness = r.toString(data["Sharpness"])
		res.NoiseReduction = r.toString(data["NoiseReduction"])
		res.ShadowTone = r.toString(data["ShadowTone"])
		res.Saturation = r.toString(data["Saturation"])
		res.ColorChromeFXBlue = r.toString(data["ColorChromeFXBlue"])
		res.ColorChromeEffect = r.toString(data["ColorChromeEffect"])
		res.GrainEffectRoughness = r.toString(data["GrainEffectRoughness"])
		res.HighlightTone = r.toString(data["HighlightTone"])
		res.LivePhotosId = r.toString(data["LivePhotosId"])
	}

	// GPS 处理
	if lat, ok := data["GPSLatitude"]; ok {
		if lng, ok := data["GPSLongitude"]; ok {
			res.GPS = map[string]float64{
				"lat": r.toFloat64(lat),
				"lng": r.toFloat64(lng),
			}
		}
	}

	return res, nil
}

func (r *Exiftool) toString(v interface{}) string {
	if v == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func (r *Exiftool) toInt(v interface{}) int {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return int(val)
	case string:
		// Handle "3000" or other numeric strings
		i, _ := strconv.Atoi(val)
		return i
	case int:
		return val
	}
	return 0
}

func (r *Exiftool) toFloat64(v interface{}) float64 {
	if v == nil {
		return 0
	}
	switch val := v.(type) {
	case float64:
		return val
	case string:
		f, _ := strconv.ParseFloat(val, 64)
		return f
	case int:
		return float64(val)
	}
	return 0
}

func filepathJoin(elem ...string) string {
	return strings.Join(elem, string(os.PathSeparator))
}
