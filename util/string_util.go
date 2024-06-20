package util

import (
	"regexp"
	"strconv"
)

func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}
func IsInArray(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func getNumbersFromString(text string) []float64 {
	if text == "" {
		return nil
	}
	// 正则表达式匹配数字（包括整数和小数）
	re := regexp.MustCompile(`\d+(\.\d+)?`)
	// 提取所有匹配的数字
	numbers := re.FindAllString(text, -1)
	floats := make([]float64, len(numbers))
	// 打印提取的数字
	for index, number := range numbers {
		floats[index], _ = strconv.ParseFloat(number, 64)
	}
	return floats
}

// ExtractPlusMinusNumbers 从输入字符串中提取所有形如 +1 或 -1 的数字部分
func ExtractPlusMinusNumbers(input string) string {
	// 定义正则表达式
	re := regexp.MustCompile(`[-+]\d+`)
	// 查找匹配项
	matches := re.FindAllString(input, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
}
