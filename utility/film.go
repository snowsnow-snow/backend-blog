package utility

import "strings"

// 定义一个映射，将富士胶片模拟的英文名称映射为中文名称
var filmModeMap = map[string]string{
	"Provia/Standard": "Provia/标准",
	"Velvia/Vivid":    "Velvia/生动",
	"Astia/Soft":      "Astia/柔和",
	"Classic Chrome":  "经典Chrome",
	"Pro Neg. Hi":     "Pro Neg. Hi",
	"Pro Neg. Std":    "Pro Neg. 标准",
	"Monochrome":      "单色",
	"Monochrome+G":    "单色+G",
	"Monochrome+R":    "单色+R",
	"Monochrome+Y":    "单色+Y",
	"ACROS":           "ACROS",
	"ACROS+G":         "ACROS+G",
	"ACROS+R":         "ACROS+R",
	"ACROS+Y":         "ACROS+Y",
}

// GetChineseFilmMode 根据英文名称获取对应的中文名称
func GetChineseFilmMode(englishName string) string {
	if chineseName, ok := filmModeMap[englishName]; ok {
		return chineseName
	}
	return englishName // 或者根据实际需求返回默认值或错误提示
}

// 定义一个映射，将富士动态范围的英文名称映射为中文名称
var dynamicRangeMap = map[string]string{
	"100":      "100",
	"200":      "200",
	"400":      "400",
	"800":      "800",
	"1600":     "1600",
	"3200":     "3200",
	"6400":     "6400",
	"12800":    "12800",
	"25600":    "25600",
	"Auto":     "自动",
	"Standard": "标准",
}

// GetChineseDynamicRange 根据英文名称获取对应的中文名称
func GetChineseDynamicRange(englishName string) string {
	if chineseName, ok := dynamicRangeMap[englishName]; ok {
		return chineseName
	}
	return englishName // 或者根据实际需求返回默认值或错误提示
}

// 定义一个映射，将富士白平衡的英文名称映射为中文名称
var whiteBalanceMap = map[string]string{
	"Auto":            "自动",
	"Daylight":        "日光",
	"Cloudy":          "多云",
	"Incandescent":    "白炽灯",
	"Fluorescent":     "荧光灯",
	"Underwater Auto": "水下自动",
	"Standard":        "标准",
}

// GetChineseWhiteBalance 根据英文名称获取对应的中文名称
func GetChineseWhiteBalance(englishName string) string {
	if chineseName, ok := whiteBalanceMap[englishName]; ok {
		return chineseName
	}
	return englishName // 或者根据实际需求返回默认值或错误提示
}

// 定义一个映射，将富士锐度的英文名称映射为中文名称
var genericDescriptionMap = map[string]string{
	"Strong": "强",
	"Hard":   "强",
	"Normal": "正常",
	"Soft":   "弱",
	"Off":    "关闭",
}

// GetChineseGenericDescriptionMap 根据英文名称获取对应的中文名称
func GetChineseGenericDescriptionMap(englishName string) string {
	if chineseName, ok := genericDescriptionMap[englishName]; ok {
		return chineseName
	}
	return englishName // 或者根据实际需求返回默认值或错误提示
}

func GetWhiteBalanceFineTuneFormat(whiteBalanceFineTune string) string {
	if whiteBalanceFineTune == "" {
		return ""
	}
	whiteBalanceFineTune = strings.ReplaceAll(whiteBalanceFineTune, "ed ", "")
	whiteBalanceFineTune = strings.ReplaceAll(whiteBalanceFineTune, ",", "")
	whiteBalanceFineTune = strings.ReplaceAll(whiteBalanceFineTune, "lue ", "")
	return whiteBalanceFineTune
}

func GetNumericAndCharParam(shadowTone string) string {
	if shadowTone == "" {
		return ""
	}
	if strings.Contains(shadowTone, "hard") || strings.Contains(shadowTone, "high") {
		return "H" + ExtractPlusMinusNumbers(shadowTone)
	}
	if strings.Contains(shadowTone, "soft") ||
		strings.Contains(shadowTone, "soft") ||
		strings.Contains(shadowTone, "strong") {
		return "S" + ExtractPlusMinusNumbers(shadowTone)
	}
	if strings.Contains(shadowTone, "normal") {
		return "0"
	}
	return "L" + ExtractPlusMinusNumbers(shadowTone)
}

func GetNumeric(str string) string {
	if str == "" {
		return ""
	}
	if strings.Contains(str, "normal") {
		return "0"
	}
	return ExtractPlusMinusNumbers(str)
}
