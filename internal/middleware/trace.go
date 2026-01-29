package middleware

import (
	"backend-blog/internal/pkg"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TraceKey 定义 Context 中的 Key 类型

func TraceMiddleware(c *fiber.Ctx) error {
	traceID := c.Get("X-Trace-ID")
	if traceID == "" {
		traceID = uuid.New().String()
	}
	ctx := context.WithValue(c.UserContext(), pkg.TraceKey, traceID)
	ctx = context.WithValue(ctx, pkg.IPKey, c.IP())
	c.SetUserContext(ctx)
	c.Set("X-Trace-ID", traceID)
	return c.Next()
}
