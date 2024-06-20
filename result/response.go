package result

import (
	constant "backend-blog"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Response struct {
	Status int     // 错误码
	Result *Result // 错误码
}
type Result struct {
	Code int         `json:"code"` // 错误码
	Msg  string      `json:"msg"`  // 错误描述
	Data interface{} `json:"data"` // 返回数据
}

const (
	successCode          = 0
	successMsg           = "ok"
	failCode             = 100
	failMsg              = "fail"
	errorCode            = 500
	errorMsg             = "error"
	SignatureInvalidCode = 401
	SignatureInvalidMsg  = "Signature invalid"
)

func successCommit(c *fiber.Ctx) {
	if transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB); transactionDB != nil {
		transactionDB.Commit()
	}
}
func Success(c *fiber.Ctx) error {
	successCommit(c)
	return NewResult(c, &Response{
		Status: fiber.StatusOK,
		Result: BuildSuccessResult(),
	})
}
func SuccessWithMsg(c *fiber.Ctx, msg string) error {
	successCommit(c)
	return NewResult(c, &Response{
		Status: fiber.StatusOK,
		Result: &Result{
			successCode,
			msg,
			nil,
		},
	})
}
func SuccessData(c *fiber.Ctx, data any) error {
	successCommit(c)
	return NewResult(c, &Response{
		Status: fiber.StatusOK,
		Result: &Result{
			successCode,
			successMsg,
			data,
		},
	})
}
func Fail(c *fiber.Ctx) error {
	return NewResult(c, &Response{
		Status: successCode,
		Result: BuildFailResult(),
	})
}
func FailWithMsg(c *fiber.Ctx, msg string) error {
	return NewResult(c, &Response{
		Status: successCode,
		Result: &Result{
			failCode,
			msg,
			nil,
		},
	})
}
func FailResult(c *fiber.Ctx, result *Result) error {
	return NewResult(c, &Response{
		Status: successCode,
		Result: result,
	})
}
func Error(c *fiber.Ctx) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if transactionDB != nil {
		transactionDB.Rollback()
	}
	return NewResult(c, &Response{
		Status: fiber.StatusInternalServerError,
		Result: BuildErrorResult(),
	})
}
func ErrorWithMsg(c *fiber.Ctx, msg string) error {
	transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB)
	if transactionDB != nil {
		transactionDB.Rollback()
	}
	return NewResult(c, &Response{
		Status: fiber.StatusInternalServerError,
		Result: &Result{
			errorCode,
			msg,
			nil,
		},
	})
}
func ErrorResult(c *fiber.Ctx, result *Result) error {
	if transactionDB := c.Locals(constant.Local.TransactionDB).(*gorm.DB); transactionDB != nil {
		transactionDB.Rollback()
	}
	return NewResult(c, &Response{
		Status: fiber.StatusInternalServerError,
		Result: result,
	})
}

func SignatureInvalidResult(c *fiber.Ctx, result *Result) error {
	return NewResult(c, &Response{
		Status: fiber.StatusUnauthorized,
		Result: result,
	})
}

func NewResult(c *fiber.Ctx, response *Response) error {
	return c.Status(response.Status).JSON(response.Result)
}

//
//// WithMsg 自定义响应信息
//func (res *Response) WithMsg(message string) Response {
//	return Response{
//		Status: res.Status,
//		Result: Result{
//			Code: res.Result.Code,
//			Msg:  message,
//			Data: res.Result.Data,
//		},
//	}
//}

// WithData 追加响应数据
func (res *Response) WithData(data interface{}) Response {
	return Response{
		Status: res.Status,
		Result: &Result{
			Code: res.Result.Code,
			Msg:  res.Result.Msg,
			Data: data,
		},
	}
}

// ToString 返回 JSON 格式的错误详情
func (res *Response) ToString() string {
	err := &struct {
		Code int         `json:"code"`
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}{
		Code: res.Result.Code,
		Msg:  res.Result.Msg,
		Data: res.Result.Data,
	}
	raw, _ := json.Marshal(err)
	return string(raw)
}

// 构造函数
func response(status int, code int, msg string, data any) *Response {
	return &Response{
		Status: status,
		Result: &Result{
			Code: code,
			Msg:  msg,
			Data: data,
		},
	}
}
func result(code int, msg string, data any) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
func BuildSuccessResult() *Result {
	return &Result{
		Code: successCode,
		Msg:  successMsg,
		Data: nil,
	}
}
func BuildSuccessResultWithMsg(msg string) *Result {
	return &Result{
		Code: successCode,
		Msg:  msg,
		Data: nil,
	}
}
func BuildSuccessResultData(msg string) *Result {
	return &Result{
		Code: successCode,
		Msg:  msg,
		Data: nil,
	}
}
func BuildFailResult() *Result {
	return &Result{
		Code: failCode,
		Msg:  failMsg,
		Data: nil,
	}
}
func BuildFailResultWithMsg(msg string) *Result {
	return &Result{
		Code: failCode,
		Msg:  msg,
		Data: nil,
	}
}
func BuildErrorResult() *Result {
	return &Result{
		Code: errorCode,
		Msg:  errorMsg,
		Data: nil,
	}
}
func BuildErrorResultWithMsg(msg string) *Result {
	return &Result{
		Code: errorCode,
		Msg:  msg,
		Data: nil,
	}
}
