package util

import (
	"github.com/dlclark/regexp2"
)

var (
	// VerifyExpUsername 要求用户名只包含字母、数字、下划线和连字符，且长度在3到20个字符之间
	VerifyExpUsername = `^[a-zA-Z0-9_-]{5,15}$`
	// VerifyExpPassword 要求密码长度大于8个字符，由大写字母、小写字母、数字、符号中的3种及3种以上组成
	VerifyExpPassword = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@._$!%*?&]{9,}$`
)

func CheckUsername(username string) (bool, error) {
	// 使用正则表达式检查用户名是否符合规则
	return Match(username, VerifyExpUsername)
}
func CheckPassword(password string) (bool, error) {
	// 使用正则表达式检查密码是否符合安全规则
	return Match(password, VerifyExpPassword)
}

func Match(str string, exp string) (bool, error) {
	regex := regexp2.MustCompile(exp, 0)
	match, err := regex.FindStringMatch(str)
	if err != nil {
		return false, err
	}
	if match != nil {
		println(match.String())
		return false, nil
	}
	return true, nil
}
