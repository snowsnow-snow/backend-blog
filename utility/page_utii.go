package utility

import "github.com/gofiber/fiber/v2"

func GetPageParam(c *fiber.Ctx) (int, int) {
	limit := c.QueryInt("limit")
	offset := c.QueryInt("offset")
	if offset == 0 {
		offset = -1
	}
	return limit, offset
}
