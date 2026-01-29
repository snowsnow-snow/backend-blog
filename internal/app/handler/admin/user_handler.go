package admin

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/middleware"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPwdReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "参数格式错误")
	}
	username := middleware.GetUsername(c.UserContext())
	err := h.userService.ResetPassword(c.UserContext(), username, req.OldPassword, req.NewPassword)
	if err != nil {
		return response.Fail(c, err.Error())
	}
	return response.Success(c)
}
