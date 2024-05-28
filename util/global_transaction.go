package util

import (
	"backend-blog"
	"backend-blog/logger"
	"backend-blog/result"
	"github.com/gofiber/fiber/v2"
)

func TransactionBegin(c *fiber.Ctx) error {
	transactionDB := DB.Begin()
	if transactionDB.Error != nil {
		logger.Error.Printf("transaction begin fail, msg: %+v", transactionDB.Error)
		return result.ErrorWithMsg(c, result.Err.Error())
	}
	c.Locals(constant.Local.TransactionDB, transactionDB)
	return c.Next()
}
