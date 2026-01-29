package common

import (
	"backend-blog/internal/model"
	"backend-blog/internal/model/entity"
	"backend-blog/internal/pkg"
	"backend-blog/utility"

	"github.com/gofiber/fiber/v2"
)

func GetCurrUsername(c *fiber.Ctx) string {
	locals := c.Locals("user")
	if locals == nil {
		return ""
	}
	return c.Locals("user").(string)
}
func CreateInit(c *fiber.Ctx, info *entity.BaseInfo) {
	if info.ID == 0 {
		info.ID = utility.GenID()
	}
	info.CreateUser = GetCurrUsername(c)
	info.CreatedTime = model.Now()
	info.PublishIp = c.IP()

	// 从 Context 中提取 Trace ID
	if tid, ok := c.UserContext().Value(pkg.TraceKey).(string); ok {
		info.TraceId = tid
	}
}
func UpdateInit(c *fiber.Ctx, info *entity.BaseInfo) {
	info.UpdateUser = GetCurrUsername(c)
	info.UpdatedTime = model.Now()
}
