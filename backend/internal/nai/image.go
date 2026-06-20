package nai

import (
	"encoding/base64"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

// DecodeDataURI 解码 data:image/png;base64,xxxx 或裸 base64，返回 mime 和原始字节。
func DecodeDataURI(s string) (mime string, data []byte, err error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "data:") {
		commaIdx := strings.Index(s, ",")
		if commaIdx < 0 {
			return "", nil, fmt.Errorf("invalid data URI: missing comma")
		}
		meta := s[:commaIdx]
		payload := s[commaIdx+1:]
		mime = "image/png"
		isBase64 := false
		for _, p := range strings.Split(meta, ";") {
			p = strings.TrimSpace(p)
			if p == "base64" {
				isBase64 = true
			} else if strings.HasPrefix(p, "image/") {
				mime = p
			}
		}
		if !isBase64 {
			return "", nil, fmt.Errorf("unsupported data URI: not base64")
		}
		data, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return "", nil, fmt.Errorf("decode base64: %w", err)
		}
		return mime, data, nil
	}
	// 裸 base64
	data, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", nil, fmt.Errorf("decode base64: %w", err)
	}
	return DetectMIME(data), data, nil
}

// DetectMIME 根据魔数探测图片 mime。
func DetectMIME(data []byte) string {
	if len(data) >= 8 && data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
		return "image/png"
	}
	if len(data) >= 3 && data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "image/jpeg"
	}
	if len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WEBP" {
		return "image/webp"
	}
	if len(data) >= 6 && (string(data[0:6]) == "GIF87a" || string(data[0:6]) == "GIF89a") {
		return "image/gif"
	}
	return "image/png"
}

// DecodeImageMeta 解码图片头部获取宽高。不支持 webp 时返回错误，调用方应忽略错误跳过尺寸校验。
func DecodeImageMeta(data []byte) (width, height int, err error) {
	cfg, _, err := image.DecodeConfig(bytesReader(data))
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}
