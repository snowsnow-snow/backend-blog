package admin

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/model/dto"
	"backend-blog/internal/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) Create(c *fiber.Ctx) error {
	var req dto.CategoryReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "参数格式错误")
	}

	if err := h.categoryService.Create(c.UserContext(), req); err != nil {
		return response.Error(c, err, "创建分类失败")
	}

	return response.Success(c)
}

func (h *CategoryHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Fail(c, "无效的分类ID")
	}

	if err := h.categoryService.Delete(c.UserContext(), id); err != nil {
		return response.Error(c, err, "删除分类失败")
	}

	return response.Success(c)
}

func (h *CategoryHandler) List(c *fiber.Ctx) error {
	result, err := h.categoryService.List(c.UserContext())
	if err != nil {
		return response.Error(c, err, "获取分类列表失败")
	}

	return response.Success(c, result)
}
