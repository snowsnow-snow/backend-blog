package response

import (
	"log/slog"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
)

// Result 业务响应结构
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"` // 为空时不显示
}

const (
	CodeSuccess = 0
	CodeFail    = 100
	CodeError   = 500
)

// --- 统一核心处理器 ---
func send(c *fiber.Ctx, httpStatus int, code int, msg string, data interface{}) error {
	return c.Status(httpStatus).JSON(Result{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// --- 成功系列 ---

func Success(c *fiber.Ctx, data ...interface{}) error {
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	return send(c, fiber.StatusOK, CodeSuccess, "ok", d)
}

// --- 失败系列 (业务逻辑错误，如密码错误) ---

func Fail(c *fiber.Ctx, msg string) error {
	return send(c, fiber.StatusOK, CodeFail, msg, nil)
}

// --- 异常系列 (系统崩溃，由中间件触发 Rollback) ---

func Error(c *fiber.Ctx, err error, msg ...string) error {
	// 此时的 err 包含了底层一路 wrap 上来的所有信息
	slog.Log(c.UserContext(), slog.LevelError, "Request Failed",
		slog.Any("error", err),
		slog.String("path", c.Path()),
		slog.String("stack", string(debug.Stack())), // 只有这里记录堆栈
	)
	m := "internal server error"
	if len(msg) > 0 {
		m = msg[0]
	}
	// 注意：这里不再进行 Rollback，回滚交给 Transaction 中间件
	return send(c, fiber.StatusInternalServerError, CodeError, m, nil)
}

func Unauthorized(c *fiber.Ctx, msg string) error {
	return send(c, fiber.StatusUnauthorized, 401, msg, nil)
}
