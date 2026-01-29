package utility

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// 注册自定义功能：使错误信息中的字段名显示为 JSON Tag 的名字
	// 例如：错误显示 "username" 而不是 "Username"
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate 统一校验方法
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err != nil {
		// 如果有校验错误，提取出第一个错误并返回可读信息
		var errs validator.ValidationErrors
		if errors.As(err, &errs) {
			for _, f := range errs {
				return fmt.Errorf("字段 [%s] 校验失败: 规则为 %s %s", f.Field(), f.Tag(), f.Param())
			}
		}
		return err
	}
	return nil
}
