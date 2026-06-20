package handlers

import (
	"encoding/json"

	"naiimage/backend/internal/models"
	"naiimage/backend/internal/nai"
	"naiimage/backend/internal/store"

	"github.com/gofiber/fiber/v2"
)

// GenerateHandler POST /api/generate
//
// 接收前端结构化绘图请求，校验 -> 组装上游 chat/completions -> 调用 -> 解析 -> 存图 -> 存历史 -> 返回。
func (h *Handler) GenerateHandler(c *fiber.Ctx) error {
	var req models.DrawRequest
	if err := c.BodyParser(&req); err != nil {
		return fiberError(c, 400, "请求体解析失败: "+err.Error())
	}

	// 1. 读取配置（base_url / api_key / 默认 model）
	s, err := store.GetSettings(h.DB, h.Defaults)
	if err != nil {
		return fiberError(c, 500, "读取配置失败: "+err.Error())
	}
	if s.UpstreamBaseURL == "" || s.UpstreamAPIKey == "" {
		return fiberError(c, 400, "未配置上游 base_url 或 api_key，请先在设置中配置")
	}

	if req.Model == "" {
		req.Model = s.DefaultModel
	}
	if req.Model == "" {
		return fiberError(c, 400, "未指定 model 且未配置 default_model")
	}

	// 2. 校验
	if err := nai.Validate(&req); err != nil {
		if ve, ok := err.(*nai.ValidationError); ok {
			return fiberError(c, ve.Status, ve.Error())
		}
		return fiberError(c, 400, err.Error())
	}

	// 3. 创建历史记录（running 状态）
	paramsJSON, _ := json.Marshal(req)
	task, err := store.CreateTask(h.DB, req.Model, paramsJSON)
	if err != nil {
		return fiberError(c, 500, "创建历史记录失败: "+err.Error())
	}

	// 4. 组装上游请求
	chatReq, err := nai.BuildChatRequest(&req)
	if err != nil {
		_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
			Status: "error",
			Error:  "组装请求失败: " + err.Error(),
		})
		return fiberError(c, 500, "组装请求失败: "+err.Error())
	}

	// 5. 调用上游
	result, err := h.Client.ChatCompletions(c.Context(), s.UpstreamBaseURL, s.UpstreamAPIKey, chatReq)
	if err != nil {
		errMsg := err.Error()
		_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
			Status: "error",
			Error:  errMsg,
		})
		return upstreamFiberError(c, err)
	}

	// 6. 解析响应内容
	parsed := nai.ParseContent(result.Content)
	if len(parsed.Images) == 0 {
		errMsg := "上游返回内容中未找到图片"
		_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
			Status:     "error",
			Error:      errMsg,
			UpstreamID: result.UpstreamID,
			Usage:      usageToRaw(result.Usage),
		})
		return fiberError(c, 502, errMsg)
	}

	// 7. 存图（取第一张，当前每次只返回 1 张）
	dataURI := parsed.Images[0]
	mime, imgBytes, decErr := nai.DecodeDataURI(dataURI)
	if decErr != nil {
		errMsg := "解析图片 data URI 失败: " + decErr.Error()
		_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
			Status:     "error",
			Error:      errMsg,
			UpstreamID: result.UpstreamID,
			Usage:      usageToRaw(result.Usage),
		})
		return fiberError(c, 502, errMsg)
	}

	width, height := 0, 0
	if w, h, e := nai.DecodeImageMeta(imgBytes); e == nil {
		width, height = w, h
	}

	img, err := store.SaveImage(h.DB, imgBytes, mime, width, height, "generated")
	if err != nil {
		errMsg := "存储图片失败: " + err.Error()
		_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
			Status:     "error",
			Error:      errMsg,
			UpstreamID: result.UpstreamID,
			Usage:      usageToRaw(result.Usage),
		})
		return fiberError(c, 500, errMsg)
	}

	// 8. 序列化 seeds / vibe_cache_ids / usage
	seedsRaw := seedsToRaw(parsed.Seeds)
	vibeRaw := vibeToRaw(parsed.VibeCacheIDs)
	usageRaw := usageToRaw(result.Usage)

	// 9. 更新历史记录为 done
	_ = store.FinishTask(h.DB, task.ID, &store.TaskResult{
		Status:        "done",
		OutputImageID: img.ID,
		Seeds:         seedsRaw,
		VibeCacheIDs:  vibeRaw,
		Usage:         usageRaw,
		UpstreamID:    result.UpstreamID,
	})

	// 10. 构造响应
	resp := models.GenerateResponse{
		TaskID: task.ID,
		Model:  result.Model,
		Status: "done",
		Image: &models.ImageRef{
			ID:     img.ID,
			URL:    "/api/images/" + img.ID,
			Mime:   img.Mime,
			Width:  img.Width,
			Height: img.Height,
		},
		Seeds:        parsed.Seeds,
		VibeCacheIDs: parsed.VibeCacheIDs,
		Usage:        usageRaw,
		UpstreamID:   result.UpstreamID,
		CreatedAt:    task.CreatedAt,
	}

	// 读取 finished_at 计算 elapsed
	if t, e := store.GetTask(h.DB, task.ID); e == nil {
		resp.FinishedAt = t.FinishedAt
		if t.FinishedAt > 0 {
			resp.ElapsedMs = (t.FinishedAt - t.CreatedAt) * 1000
		}
	}

	return c.JSON(resp)
}

// usageToRaw 直接透传上游返回的 usage JSON（已经是 RawMessage）。
func usageToRaw(u json.RawMessage) json.RawMessage {
	if len(u) == 0 {
		return nil
	}
	return u
}

// seedsToRaw 把 seeds 切片序列化成 json.RawMessage。
func seedsToRaw(seeds []int64) json.RawMessage {
	if len(seeds) == 0 {
		return nil
	}
	b, err := json.Marshal(seeds)
	if err != nil {
		return nil
	}
	return b
}

// vibeToRaw 把 vibe_cache_ids 切片序列化成 json.RawMessage。
func vibeToRaw(ids []models.VibeCacheID) json.RawMessage {
	if len(ids) == 0 {
		return nil
	}
	b, err := json.Marshal(ids)
	if err != nil {
		return nil
	}
	return b
}
