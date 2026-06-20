package nai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"naiimage/backend/internal/models"
)

// Client 是上游 NewAPI 客户端。
type Client struct {
	httpClient *http.Client
}

// NewClient 创建上游客户端。
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: timeout},
	}
}

// normalizeBaseURL 去掉末尾斜杠。
func normalizeBaseURL(url string) string {
	return strings.TrimRight(url, "/")
}

// ModelsResult 是 GetModels 的结果。
type ModelsResult struct {
	Models []models.ModelInfo
	Raw    []byte
}

// GetModels 调用上游 /v1/models 获取模型列表。
func (c *Client) GetModels(ctx context.Context, baseURL, apiKey string) (*ModelsResult, error) {
	baseURL = normalizeBaseURL(baseURL)
	url := baseURL + "/v1/models"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build models request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("models request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, fmt.Errorf("read models response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ParseUpstreamError(resp.StatusCode, body)
	}

	var mr models.ModelsResponse
	if err := json.Unmarshal(body, &mr); err != nil {
		return nil, fmt.Errorf("parse models response: %w", err)
	}
	return &ModelsResult{Models: mr.Data, Raw: body}, nil
}

// ChatResult 是 ChatCompletions 的结果。
type ChatResult struct {
	Content    string
	Usage      json.RawMessage
	UpstreamID string
	Model      string
	Raw        []byte
}

// ChatCompletions 调用上游 /v1/chat/completions 绘图。
func (c *Client) ChatCompletions(ctx context.Context, baseURL, apiKey string, cr *chatRequest) (*ChatResult, error) {
	baseURL = normalizeBaseURL(baseURL)
	url := baseURL + "/v1/chat/completions"

	payload, err := json.Marshal(cr)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build chat request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chat request failed: %w", err)
	}
	defer resp.Body.Close()

	// 绘图返回可能很大（base64 图片），不限制读取大小
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read chat response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ParseUpstreamError(resp.StatusCode, body)
	}

	var chatResp chatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("parse chat response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, &UpstreamError{Status: resp.StatusCode, Message: "上游返回 choices 为空"}
	}

	content := ""
	if chatResp.Choices[0].Message.Content != "" {
		content = chatResp.Choices[0].Message.Content
	}

	var usageRaw json.RawMessage
	if chatResp.Usage != nil {
		if b, err := json.Marshal(chatResp.Usage); err == nil {
			usageRaw = b
		}
	}

	return &ChatResult{
		Content:    content,
		Usage:      usageRaw,
		UpstreamID: chatResp.ID,
		Model:      chatResp.Model,
		Raw:        body,
	}, nil
}
