package client

import (
	"backend-blog/internal/app/services"
	"backend-blog/internal/pkg/response"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type MediaHandler struct {
	service *services.MediaService
}

func NewMediaHandler() *MediaHandler {
	return &MediaHandler{
		service: services.NewMediaService(),
	}
}

func (h *MediaHandler) ListByPost(c *fiber.Ctx) error {
	postIDStr := c.Params("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Post ID")
	}

	assets, err := h.service.GetByPostIDClient(c.UserContext(), postID)
	if err != nil {
		return response.Error(c, err, "Failed to fetch media assets: "+err.Error())
	}

	return response.Success(c, assets)
}

func (h *MediaHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Media ID")
	}

	asset, err := h.service.GetByID(c.UserContext(), id)
	if err != nil {
		return response.Error(c, err, "Media not found")
	}

	// Serve compressed version by default if it exists
	targetPath := asset.ThumbnailPath
	if targetPath == "" {
		targetPath = asset.FilePath
	}

	return c.SendFile(targetPath)
}
