package handlers

import (
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
)

// SettingsResponse 返回给前端的配置，api_key 脱敏（只返回是否已设置）。
type SettingsResponse struct {
	UpstreamBaseURL string `json:"upstream_base_url"`
	HasAPIKey       bool   `json:"has_api_key"`
	DefaultModel    string `json:"default_model"`
	RequestTimeout  int    `json:"request_timeout_seconds"`
}

// SettingsUpdateRequest 前端更新配置的请求体。api_key 为空表示不修改。
type SettingsUpdateRequest struct {
	UpstreamBaseURL string `json:"upstream_base_url"`
	UpstreamAPIKey  string `json:"upstream_api_key"`
	DefaultModel    string `json:"default_model"`
	RequestTimeout  int    `json:"request_timeout_seconds"`
}

// GetSettingsHandler GET /api/settings
func (h *Handler) GetSettingsHandler(c *fiber.Ctx) error {
	s, err := store.GetSettings(h.DB, h.Defaults)
	if err != nil {
		return fiberError(c, 500, "读取配置失败: "+err.Error())
	}
	return c.JSON(SettingsResponse{
		UpstreamBaseURL: s.UpstreamBaseURL,
		HasAPIKey:       s.UpstreamAPIKey != "",
		DefaultModel:    s.DefaultModel,
		RequestTimeout:  s.RequestTimeout,
	})
}

// UpdateSettingsHandler PUT /api/settings
func (h *Handler) UpdateSettingsHandler(c *fiber.Ctx) error {
	var req SettingsUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiberError(c, 400, "请求体解析失败: "+err.Error())
	}

	// 读取当前配置，api_key 为空时保留原值
	current, err := store.GetSettings(h.DB, h.Defaults)
	if err != nil {
		return fiberError(c, 500, "读取配置失败: "+err.Error())
	}

	next := store.Settings{
		UpstreamBaseURL: req.UpstreamBaseURL,
		UpstreamAPIKey:  req.UpstreamAPIKey,
		DefaultModel:    req.DefaultModel,
		RequestTimeout:  req.RequestTimeout,
	}
	if next.UpstreamBaseURL == "" {
		next.UpstreamBaseURL = current.UpstreamBaseURL
	}
	if next.UpstreamAPIKey == "" {
		next.UpstreamAPIKey = current.UpstreamAPIKey
	}
	if next.DefaultModel == "" {
		next.DefaultModel = current.DefaultModel
	}
	if next.RequestTimeout <= 0 {
		next.RequestTimeout = current.RequestTimeout
	}

	if err := store.UpdateSettings(h.DB, next); err != nil {
		return fiberError(c, 500, "保存配置失败: "+err.Error())
	}

	return c.JSON(SettingsResponse{
		UpstreamBaseURL: next.UpstreamBaseURL,
		HasAPIKey:       next.UpstreamAPIKey != "",
		DefaultModel:    next.DefaultModel,
		RequestTimeout:  next.RequestTimeout,
	})
}
