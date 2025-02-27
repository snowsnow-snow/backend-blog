package common

import (
	"backend-blog/internal/model"
	"backend-blog/internal/model/entity"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetCurrUsername(c *fiber.Ctx) string {
	locals := c.Locals("user")
	if locals == nil {
		return ""
	}
	return c.Locals("user").(string)
}
func CreateInit(c *fiber.Ctx, info *entity.BaseInfo) {
	if info.ID == "" {
		info.ID = uuid.NewString()
	}
	info.CreateUser = GetCurrUsername(c)
	info.CreatedTime = model.Now()
	info.PublishIp = c.IP()
	info.RequestId = c.Locals("Request-ID").(string)
}
func UpdateInit(c *fiber.Ctx, info *entity.BaseInfo) {
	info.UpdateUser = GetCurrUsername(c)
	info.UpdatedTime = model.Now()
}
