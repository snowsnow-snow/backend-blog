package admin

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/model/dto"
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

func (h *PostHandler) Create(c *fiber.Ctx) error {
	var req dto.CreatePostReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "Invalid request body")
	}

	if req.Title == "" {
		return response.Fail(c, "Title is required")
	}

	if err := h.service.CreatePost(c.UserContext(), req); err != nil {
		return response.Error(c, err, "Failed to create post: "+err.Error())
	}

	return response.Success(c, nil)
}

func (h *PostHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Post ID")
	}

	var req dto.UpdatePostReq
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, "Invalid request body")
	}
	req.ID = id

	if err := h.service.UpdatePost(c.UserContext(), req); err != nil {
		return response.Error(c, err, "Failed to update post: "+err.Error())
	}

	return response.Success(c, nil)
}

func (h *PostHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Post ID")
	}

	if err := h.service.DeletePost(c.UserContext(), id); err != nil {
		return response.Error(c, err, "Failed to delete post: "+err.Error())
	}

	return response.Success(c, nil)
}

func (h *PostHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Post ID")
	}

	post, err := h.service.GetPostAdmin(c.UserContext(), id)
	if err != nil {
		return response.Error(c, err, "Failed to get post: "+err.Error())
	}

	return response.Success(c, post)
}

func (h *PostHandler) List(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)

	posts, total, err := h.service.ListPostsAdmin(c.UserContext(), page, pageSize)
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
