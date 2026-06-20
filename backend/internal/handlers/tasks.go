package handlers

import (
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
)

// TaskDTO 历史记录对外表示。
type TaskDTO struct {
	ID            string      `json:"id"`
	Model         string      `json:"model"`
	Params        interface{} `json:"params"`
	Status        string      `json:"status"`
	Error         string      `json:"error,omitempty"`
	OutputImageID string      `json:"output_image_id,omitempty"`
	ImageURL      string      `json:"image_url,omitempty"`
	Seeds         interface{} `json:"seeds,omitempty"`
	VibeCacheIDs  interface{} `json:"vibe_cache_ids,omitempty"`
	Usage         interface{} `json:"usage,omitempty"`
	UpstreamID    string      `json:"upstream_id,omitempty"`
	CreatedAt     int64       `json:"created_at"`
	FinishedAt    int64       `json:"finished_at,omitempty"`
	ElapsedMs     int64       `json:"elapsed_ms,omitempty"`
}

// toTaskDTO 把 store.Task 转成 DTO，params/seeds/usage 从 RawMessage 反序列化为 interface{}。
func toTaskDTO(t *store.Task) TaskDTO {
	dto := TaskDTO{
		ID:            t.ID,
		Model:         t.Model,
		Status:        t.Status,
		Error:         t.Error,
		OutputImageID: t.OutputImageID,
		UpstreamID:    t.UpstreamID,
		CreatedAt:     t.CreatedAt,
		FinishedAt:    t.FinishedAt,
	}
	if t.OutputImageID != "" {
		dto.ImageURL = "/api/images/" + t.OutputImageID
	}
	dto.Params = rawToInterface(t.Params)
	dto.Seeds = rawToInterface(t.Seeds)
	dto.VibeCacheIDs = rawToInterface(t.VibeCacheIDs)
	dto.Usage = rawToInterface(t.Usage)
	if t.FinishedAt > 0 {
		dto.ElapsedMs = (t.FinishedAt - t.CreatedAt) * 1000
	}
	return dto
}

// ListTasksHandler GET /api/tasks?limit=50&offset=0
func (h *Handler) ListTasksHandler(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	tasks, total, err := store.ListTasks(h.DB, limit, offset)
	if err != nil {
		return fiberError(c, 500, "查询历史失败: "+err.Error())
	}

	dtos := make([]TaskDTO, 0, len(tasks))
	for _, t := range tasks {
		dtos = append(dtos, toTaskDTO(t))
	}
	return c.JSON(fiber.Map{
		"tasks": dtos,
		"total": total,
	})
}

// GetTaskHandler GET /api/tasks/:id
func (h *Handler) GetTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	t, err := store.GetTask(h.DB, id)
	if err != nil {
		return fiberError(c, 404, "任务不存在")
	}
	return c.JSON(toTaskDTO(t))
}

// DeleteTaskHandler DELETE /api/tasks/:id
func (h *Handler) DeleteTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// 先取出关联的输出图片，删除任务后一并清理图片
	t, _ := store.GetTask(h.DB, id)
	if t != nil && t.OutputImageID != "" {
		_ = store.DeleteImage(h.DB, t.OutputImageID)
	}

	if err := store.DeleteTask(h.DB, id); err != nil {
		return fiberError(c, 500, "删除任务失败: "+err.Error())
	}
	return c.JSON(fiber.Map{"ok": true})
}

// DeleteAllTasksHandler DELETE /api/tasks —— 清空全部历史
func (h *Handler) DeleteAllTasksHandler(c *fiber.Ctx) error {
	if err := store.DeleteAllTasks(h.DB); err != nil {
		return fiberError(c, 500, "清空历史失败: "+err.Error())
	}
	return c.JSON(fiber.Map{"ok": true})
}
