package nai

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"naiimage/backend/internal/models"
)

// markdown 图片：![xxx](data:image/png;base64,...)
var imgRe = regexp.MustCompile(`!\[[^\]]*\]\((data:image/[^;)]+;base64,[^)]+)\)`)

// seeds 注释：<!-- seeds:[123456789] -->
var seedRe = regexp.MustCompile(`<!--\s*seeds:(\[.*?\])\s*-->`)

// vibe_cache_ids 注释：<!-- vibe_cache_ids:[...] -->
var vibeRe = regexp.MustCompile(`<!--\s*vibe_cache_ids:(\[.*?\])\s*-->`)

// ParsedContent 解析上游 message.content 后的结果。
type ParsedContent struct {
	Images       []string // data URI 列表
	Seeds        []int64  // seeds
	VibeCacheIDs []models.VibeCacheID
}

// ParseContent 解析上游返回的 markdown content，提取图片 data URI、seeds、vibe_cache_ids。
func ParseContent(content string) ParsedContent {
	result := ParsedContent{}

	matches := imgRe.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		if len(m) >= 2 {
			result.Images = append(result.Images, m[1])
		}
	}

	if m := seedRe.FindStringSubmatch(content); m != nil && len(m) >= 2 {
		var seeds []int64
		if err := json.Unmarshal([]byte(m[1]), &seeds); err == nil {
			result.Seeds = seeds
		}
	}

	if m := vibeRe.FindStringSubmatch(content); m != nil && len(m) >= 2 {
		var ids []models.VibeCacheID
		if err := json.Unmarshal([]byte(m[1]), &ids); err == nil {
			result.VibeCacheIDs = ids
		}
	}

	return result
}

// ParseUpstreamError 把上游非 2xx 响应体解析成 UpstreamError。
func ParseUpstreamError(statusCode int, body []byte) *UpstreamError {
	ue := &UpstreamError{Status: statusCode}

	// 尝试 OpenAI 风格 {"error":{"message":...}}
	var errWrap struct {
		Error *upstreamErrBody `json:"error"`
	}
	if err := json.Unmarshal(body, &errWrap); err == nil && errWrap.Error != nil {
		msg := errWrap.Error.Message
		if msg == "" {
			msg = fmt.Sprintf("upstream error (HTTP %d)", statusCode)
		}
		ue.Message = msg
		return ue
	}

	// 尝试纯字符串
	bodyStr := strings.TrimSpace(string(body))
	if bodyStr != "" {
		ue.Message = fmt.Sprintf("upstream error (HTTP %d): %s", statusCode, bodyStr)
		return ue
	}

	ue.Message = fmt.Sprintf("upstream error (HTTP %d)", statusCode)
	return ue
}
