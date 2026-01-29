package client

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	result, err := h.categoryService.List(c.UserContext())
	if err != nil {
		return response.Error(c, err, "获取分类列表失败")
	}

	return response.Success(c, result)
}
