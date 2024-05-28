package common

import (
	"backend-blog/config"
	"backend-blog/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
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
func GetCurrUser(c *fiber.Ctx) map[string]interface{} {
	return c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
}

func GetCurrUsername(c *fiber.Ctx) string {
	return GetCurrUser(c)["username"].(string)
}
func CreateInit(c *fiber.Ctx, info *models.BaseInfo) {
	info.ID = uuid.NewString()
	//info.CreateUser = GetCurrUsername(c)
	info.CreatedTime = models.DateTime{}.Now()
	info.PublishIp = c.IP()
	info.RequestId = c.Locals("Request-ID").(string)
}
func UpdateInit(info *models.BaseInfo) {
	//now := 	models.DateTime{}.Now()
	//info.Id = uuid.NewString()
	//info.UpdateUser = GetCurrUsername(c)
	info.UpdatedTime = models.DateTime{}.Now()
}
