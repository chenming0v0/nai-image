package nai

import (
	"encoding/json"

	"naiimage/backend/internal/models"
)

// innerPayload 是放入 messages[0].content 的真正绘图参数对象。
// 注意：不包含 model（由外层 body.model 提供）、不包含 max_tokens / stream。
type innerPayload struct {
	Prompt         string                   `json:"prompt"`
	NegativePrompt string                   `json:"negative_prompt,omitempty"`
	Size           []int                    `json:"size,omitempty"`
	Steps          *int                     `json:"steps,omitempty"`
	Scale          *float64                 `json:"scale,omitempty"`
	Sampler        string                   `json:"sampler,omitempty"`
	Seed           *int64                   `json:"seed,omitempty"`
	VarietyBoost   *bool                    `json:"variety_boost,omitempty"`
	CFGRescale     *float64                 `json:"cfg_rescale,omitempty"`
	NoiseSchedule  string                   `json:"noise_schedule,omitempty"`
	ImageFormat    string                   `json:"image_format,omitempty"`
	NSamples       *int                     `json:"n_samples,omitempty"`
	UseCoords      *bool                    `json:"use_coords,omitempty"`
	UseOrder       *bool                    `json:"use_order,omitempty"`
	Characters     []models.Character       `json:"characters,omitempty"`
	I2I            *models.I2IParams        `json:"i2i,omitempty"`
	Inpaint        *models.InpaintParams    `json:"inpaint,omitempty"`
	Controlnet     *models.ControlnetParams `json:"controlnet,omitempty"`
	CharacterRefs  []models.CharacterRef    `json:"character_references,omitempty"`
}

// chatRequest 是发给上游 /v1/chat/completions 的外层 OpenAI 兼容请求体。
type chatRequest struct {
	Model     string        `json:"model"`
	Messages  []chatMessage `json:"messages"`
	Stream    bool          `json:"stream"`
	MaxTokens int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatResponse 是上游返回的 OpenAI 兼容响应。
type chatResponse struct {
	ID      string           `json:"id"`
	Object  string           `json:"object"`
	Created int64            `json:"created"`
	Model   string           `json:"model"`
	Choices []chatChoice     `json:"choices"`
	Usage   *chatUsage       `json:"usage,omitempty"`
	Error   *upstreamErrBody `json:"error,omitempty"`
}

type chatChoice struct {
	Index        int           `json:"index"`
	Message      chatChoiceMsg `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type chatChoiceMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type upstreamErrBody struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}

// BuildChatRequest 把前端结构化 DrawRequest 组装成上游 chat/completions 请求体。
func BuildChatRequest(req *models.DrawRequest) (*chatRequest, error) {
	inner := innerPayload{
		Prompt:         req.Prompt,
		NegativePrompt: req.NegativePrompt,
		Size:           req.Size,
		Steps:          req.Steps,
		Scale:          req.Scale,
		Sampler:        req.Sampler,
		Seed:           req.Seed,
		VarietyBoost:   req.VarietyBoost,
		CFGRescale:     req.CFGRescale,
		NoiseSchedule:  req.NoiseSchedule,
		ImageFormat:    req.ImageFormat,
		NSamples:       req.NSamples,
		UseCoords:      req.UseCoords,
		UseOrder:       req.UseOrder,
		Characters:     req.Characters,
		I2I:            req.I2I,
		Inpaint:        req.Inpaint,
		Controlnet:     req.Controlnet,
		CharacterRefs:  req.CharacterRefs,
	}
	innerBytes, err := json.Marshal(inner)
	if err != nil {
		return nil, err
	}

	maxTokens := req.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 100000
	}

	stream := false
	if req.Stream != nil {
		stream = *req.Stream
	}

	return &chatRequest{
		Model:     req.Model,
		Messages:  []chatMessage{{Role: "user", Content: string(innerBytes)}},
		Stream:    stream,
		MaxTokens: maxTokens,
	}, nil
}
