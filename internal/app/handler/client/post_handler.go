package client

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type PostHandler struct {
	service *services.PostService
}

func NewPostHandler() *PostHandler {
	return &PostHandler{
		service: services.NewPostService(),
	}
}

func (h *PostHandler) Get(c *fiber.Ctx) error {
	param := c.Params("id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid post ID")
	}

	post, err := h.service.GetPostClient(c.UserContext(), id)
	if err != nil {
		return response.Error(c, err, "Failed to get post")
	}
	if post == nil {
		return response.Fail(c, "Post not found or unavailable")
	}

	return response.Success(c, post)
}

func (h *PostHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	posts, total, err := h.service.ListPostsClient(c.UserContext(), page, pageSize)
	if err != nil {
		return response.Error(c, err, "Failed to list posts: "+err.Error())
	}

	return response.Success(c, fiber.Map{
		"list":     posts,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}
