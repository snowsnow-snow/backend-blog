package common

import (
	"backend-blog/config"
	"backend-blog/models"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func init() {
	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatal("Error read config file", err)
	}
	// yaml 文件内容影射到结构体中
	err = yaml.Unmarshal(yamlFile, &config.GlobalConfig)
	if err != nil {
		log.Fatal("Error init config", err)
	}
}
func GetCurrUsername(c *fiber.Ctx) string {
	locals := c.Locals("user")
	if locals == nil {
		return ""
	}
	return c.Locals("user").(string)
}
func CreateInit(c *fiber.Ctx, info *models.BaseInfo) {
	info.ID = uuid.NewString()
	info.CreateUser = GetCurrUsername(c)
	info.CreatedTime = models.DateTime{}.Now()
	info.PublishIp = c.IP()
	info.RequestId = c.Locals("Request-ID").(string)
}
func UpdateInit(c *fiber.Ctx, info *models.BaseInfo) {
	info.UpdateUser = GetCurrUsername(c)
	info.UpdatedTime = models.DateTime{}.Now()
}
