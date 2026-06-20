package handlers

import (
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
)

// GetImageHandler GET /api/images/:id —— 返回图片二进制
func (h *Handler) GetImageHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiberError(c, 400, "缺少图片 id")
	}

	data, mime, err := store.GetImageData(h.DB, id)
	if err != nil {
		return fiberError(c, 404, "图片不存在")
	}

	c.Set("Content-Type", mime)
	c.Set("Cache-Control", "public, max-age=31536000, immutable")
	return c.Send(data)
}

// GetImageMetaHandler GET /api/images/:id/meta —— 返回图片元信息
func (h *Handler) GetImageMetaHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiberError(c, 400, "缺少图片 id")
	}

	img, err := store.GetImageMeta(h.DB, id)
	if err != nil {
		return fiberError(c, 404, "图片不存在")
	}

	return c.JSON(fiber.Map{
		"id":         img.ID,
		"mime":       img.Mime,
		"width":      img.Width,
		"height":     img.Height,
		"source":     img.Source,
		"created_at": img.CreatedAt,
	})
}
