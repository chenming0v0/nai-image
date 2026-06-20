package models

import "encoding/json"

// Settings 后端配置，对应 SQLite settings 表。
type Settings struct {
	UpstreamBaseURL string `json:"upstream_base_url"`
	UpstreamAPIKey  string `json:"upstream_api_key"`
	DefaultModel    string `json:"default_model"`
	RequestTimeout  int    `json:"request_timeout_seconds"`
}

// Character 多角色定义。
type Character struct {
	Prompt         string `json:"prompt"`
	NegativePrompt string `json:"negative_prompt,omitempty"`
	Position       string `json:"position,omitempty"`
}

// I2IParams 图生图参数。
type I2IParams struct {
	Image    string   `json:"image"`
	Strength *float64 `json:"strength,omitempty"`
	Noise    *float64 `json:"noise,omitempty"`
	Seed     *int64   `json:"seed,omitempty"`
}

// InpaintParams 局部重绘参数。
type InpaintParams struct {
	Image    string   `json:"image"`
	Mask     string   `json:"mask"`
	Strength *float64 `json:"strength,omitempty"`
	Seed     *int64   `json:"seed,omitempty"`
}

// ControlnetImage Vibe Transfer 单张参考图。
type ControlnetImage struct {
	Image         string   `json:"image,omitempty"`
	CacheID       string   `json:"cache_id,omitempty"`
	InfoExtracted *float64 `json:"info_extracted,omitempty"`
	Strength      *float64 `json:"strength,omitempty"`
}

// ControlnetParams Vibe Transfer 参数。
type ControlnetParams struct {
	Strength *float64          `json:"strength,omitempty"`
	Images   []ControlnetImage `json:"images"`
}

// CharacterRef 角色参考图。
type CharacterRef struct {
	Image    string   `json:"image"`
	Type     string   `json:"type,omitempty"`
	Fidelity *float64 `json:"fidelity,omitempty"`
	Strength *float64 `json:"strength,omitempty"`
}

// DrawRequest 前端发给后端的绘图请求（结构化、友好）。
type DrawRequest struct {
	Model          string          `json:"model"`
	Prompt         string          `json:"prompt"`
	NegativePrompt string          `json:"negative_prompt,omitempty"`
	Size           []int           `json:"size,omitempty"`
	Steps          *int            `json:"steps,omitempty"`
	Scale          *float64        `json:"scale,omitempty"`
	Sampler        string          `json:"sampler,omitempty"`
	Seed           *int64          `json:"seed,omitempty"`
	VarietyBoost   *bool           `json:"variety_boost,omitempty"`
	CFGRescale     *float64        `json:"cfg_rescale,omitempty"`
	NoiseSchedule  string          `json:"noise_schedule,omitempty"`
	ImageFormat    string          `json:"image_format,omitempty"`
	NSamples       *int            `json:"n_samples,omitempty"`
	UseCoords      *bool           `json:"use_coords,omitempty"`
	UseOrder       *bool           `json:"use_order,omitempty"`
	Characters     []Character     `json:"characters,omitempty"`
	I2I            *I2IParams      `json:"i2i,omitempty"`
	Inpaint        *InpaintParams  `json:"inpaint,omitempty"`
	Controlnet     *ControlnetParams `json:"controlnet,omitempty"`
	CharacterRefs  []CharacterRef  `json:"character_references,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Stream         *bool           `json:"stream,omitempty"`
}

// GenerateResponse 后端返回给前端的绘图结果。
type GenerateResponse struct {
	TaskID       string          `json:"task_id"`
	Model        string          `json:"model"`
	Status       string          `json:"status"`
	Image        *ImageRef       `json:"image,omitempty"`
	Seeds        []int64         `json:"seeds"`
	VibeCacheIDs []VibeCacheID   `json:"vibe_cache_ids,omitempty"`
	Usage        json.RawMessage `json:"usage,omitempty"`
	UpstreamID   string          `json:"upstream_id,omitempty"`
	Error        string          `json:"error,omitempty"`
	CreatedAt    int64           `json:"created_at"`
	FinishedAt   int64           `json:"finished_at,omitempty"`
	ElapsedMs    int64           `json:"elapsed_ms,omitempty"`
}

// VibeCacheID 对应文档 vibe_cache_ids 注释里的条目。
type VibeCacheID struct {
	Index   int    `json:"index"`
	CacheID string `json:"cache_id"`
}

// ImageRef 图片引用，url 是后端取图地址。
type ImageRef struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Mime   string `json:"mime"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// ImageMeta 图片元信息。
type ImageMeta struct {
	ID        string `json:"id"`
	Mime      string `json:"mime"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Source    string `json:"source"`
	CreatedAt int64  `json:"created_at"`
}

// TaskRecord 历史记录。
type TaskRecord struct {
	ID            string          `json:"id"`
	Model         string          `json:"model"`
	Params        json.RawMessage `json:"params"`
	Status        string          `json:"status"`
	Error         string          `json:"error,omitempty"`
	OutputImageID string          `json:"output_image_id,omitempty"`
	ImageURL      string          `json:"image_url,omitempty"`
	Seeds         json.RawMessage `json:"seeds,omitempty"`
	VibeCacheIDs  json.RawMessage `json:"vibe_cache_ids,omitempty"`
	Usage         json.RawMessage `json:"usage,omitempty"`
	UpstreamID    string          `json:"upstream_id,omitempty"`
	CreatedAt     int64           `json:"created_at"`
	FinishedAt    int64           `json:"finished_at,omitempty"`
}

// TaskListResponse 历史列表。
type TaskListResponse struct {
	Tasks []TaskRecord `json:"tasks"`
	Total int          `json:"total"`
}

// ModelInfo 上游模型列表项。
type ModelInfo struct {
	ID      string `json:"id"`
	Object  string `json:"object,omitempty"`
	Created int64  `json:"created,omitempty"`
	OwnedBy string `json:"owned_by,omitempty"`
}

// ModelsResponse 模型列表响应。
type ModelsResponse struct {
	Object string     `json:"object"`
	Data   []ModelInfo `json:"data"`
}
