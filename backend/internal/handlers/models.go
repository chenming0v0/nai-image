package handlers

import (
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
)

// GetModelsHandler GET /api/models
func (h *Handler) GetModelsHandler(c *fiber.Ctx) error {
	s, err := store.GetSettings(h.DB, h.Defaults)
	if err != nil {
		return fiberError(c, 500, "读取配置失败: "+err.Error())
	}
	if s.UpstreamBaseURL == "" || s.UpstreamAPIKey == "" {
		return fiberError(c, 400, "未配置上游 base_url 或 api_key，请先在设置中配置")
	}

	result, err := h.Client.GetModels(c.Context(), s.UpstreamBaseURL, s.UpstreamAPIKey)
	if err != nil {
		return upstreamFiberError(c, err)
	}

	// 直接透传上游原始 JSON，保证前端拿到完整结构
	c.Set("Content-Type", "application/json")
	return c.Send(result.Raw)
}
