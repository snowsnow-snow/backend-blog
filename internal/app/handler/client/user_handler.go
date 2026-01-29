package client

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/pkg/response"
	"backend-blog/utility"

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

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "参数格式错误")
	}
	res, err := h.userService.Login(c.UserContext(), req)
	if err != nil {
		return response.Fail(c, err.Error())
	}
	return response.Success(c, res)
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "格式错误")
	}
	if err := utility.Validate(req); err != nil {
		return response.Fail(c, err.Error())
	}
	if err := h.userService.Add(c.UserContext(), req); err != nil {
		return response.Fail(c, err.Error())
	}
	return response.Success(c)
}
