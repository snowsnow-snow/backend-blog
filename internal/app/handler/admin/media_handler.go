package admin

import (
	"backend-blog/config"
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

func (h *MediaHandler) Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return response.Fail(c, "Failed to get file from request")
	}

	uploadDir := "public/uploads"
	if config.GlobalConfig.File.Path.System != "" {
		uploadDir = config.GlobalConfig.File.Path.System + config.GlobalConfig.File.Path.Resource
	}

	var postID *int64
	if pidStr := c.FormValue("postId"); pidStr != "" {
		if id, err := strconv.ParseInt(pidStr, 10, 64); err == nil {
			postID = &id
		}
	}

	livePhotoID := c.FormValue("livePhotoId")

	asset, err := h.service.UploadAndProcess(c.UserContext(), file, uploadDir, postID, livePhotoID)
	if err != nil {
		return response.Error(c, err, "Failed to process upload: "+err.Error())
	}

	return response.Success(c, asset)
}

func (h *MediaHandler) ListByPost(c *fiber.Ctx) error {
	postIDStr := c.Params("id")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		return response.Fail(c, "Invalid Post ID")
	}

	assets, err := h.service.GetByPostIDAdmin(c.UserContext(), postID)
	if err != nil {
		return response.Error(c, err, "Failed to fetch media assets: "+err.Error())
	}

	return response.Success(c, assets)
}
