package nai

import (
	"regexp"
	"strconv"
	"strings"

	"naiimage/backend/internal/models"
)

// 允许的尺寸上限（文档 §6）
const (
	maxWidth  = 1216
	maxHeight = 1216
	maxPixels = 1024 * 1024
)

// 允许的采样器
var allowedSamplers = map[string]bool{
	"k_euler": true, "k_euler_ancestral": true, "k_dpm_2": true,
	"k_dpm_2_ancestral": true, "k_dpmpp_2m": true,
	"k_dpmpp_2s_ancestral": true, "k_dpmpp_sde": true, "ddim": true,
}

// 允许的 noise schedule
var allowedNoiseSchedules = map[string]bool{
	"karras": true, "exponential": true, "polyexponential": true,
}

// 允许的图片格式
var allowedImageFormats = map[string]bool{
	"png": true, "webp": true,
}

// 允许的 character_references type
var allowedCharRefTypes = map[string]bool{
	"character": true, "style": true, "character&style": true,
}

// CJK / 假名 / 韩文 / 全角符号检测
var cjkRe = regexp.MustCompile(`[\x{3040}-\x{30FF}\x{3400}-\x{4DBF}\x{4E00}-\x{9FFF}\x{AC00}-\x{D7AF}\x{FF00}-\x{FFEF}]`)

// position 网格 [A-E][1-5]
var positionRe = regexp.MustCompile(`^[A-E][1-5]$`)

// Validate 校验绘图请求。返回 nil 表示通过。
func Validate(req *models.DrawRequest) error {
	if req == nil {
		return newValidationError("", "request is nil")
	}

	// prompt 必填
	if strings.TrimSpace(req.Prompt) == "" {
		return newValidationError("prompt", "prompt 不能为空")
	}
	if err := checkNoCJK("prompt", req.Prompt); err != nil {
		return err
	}
	if err := checkNoCJK("negative_prompt", req.NegativePrompt); err != nil {
		return err
	}

	// model 必填
	if strings.TrimSpace(req.Model) == "" {
		return newValidationError("model", "model 不能为空")
	}

	// size 校验
	var sizeW, sizeH int
	if len(req.Size) > 0 {
		if len(req.Size) != 2 {
			return newValidationError("size", "size 必须是 [width, height] 整数数组")
		}
		sizeW = req.Size[0]
		sizeH = req.Size[1]
		if err := validateSize(sizeW, sizeH); err != nil {
			return err
		}
	}

	// steps
	if req.Steps != nil {
		if *req.Steps < 1 || *req.Steps > 28 {
			return newValidationError("steps", "steps 范围 1~28")
		}
	}

	// n_samples
	if req.NSamples != nil && *req.NSamples != 1 {
		return newValidationError("n_samples", "当前只允许 1")
	}

	// sampler
	if req.Sampler != "" && !allowedSamplers[req.Sampler] {
		return newValidationError("sampler", "不支持的采样器: "+req.Sampler)
	}

	// noise_schedule
	if req.NoiseSchedule != "" && !allowedNoiseSchedules[req.NoiseSchedule] {
		return newValidationError("noise_schedule", "不支持的 noise_schedule: "+req.NoiseSchedule)
	}

	// image_format
	if req.ImageFormat != "" && !allowedImageFormats[req.ImageFormat] {
		return newValidationError("image_format", "只支持 png 或 webp")
	}

	// scale
	if req.Scale != nil {
		if *req.Scale < 0 || *req.Scale > 10 {
			return newValidationError("scale", "scale 范围 0~10")
		}
	}

	// cfg_rescale
	if req.CFGRescale != nil {
		if *req.CFGRescale < 0 || *req.CFGRescale > 1 {
			return newValidationError("cfg_rescale", "cfg_rescale 范围 0~1")
		}
	}

	// characters
	for i, c := range req.Characters {
		prefix := "characters[" + strconv.Itoa(i) + "]"
		if strings.TrimSpace(c.Prompt) == "" {
			return newValidationError(prefix+".prompt", "角色 prompt 不能为空")
		}
		if err := checkNoCJK(prefix+".prompt", c.Prompt); err != nil {
			return err
		}
		if err := checkNoCJK(prefix+".negative_prompt", c.NegativePrompt); err != nil {
			return err
		}
		if c.Position != "" && !positionRe.MatchString(c.Position) {
			return newValidationError(prefix+".position", "position 格式为 [A-E][1-5]，如 C3")
		}
	}

	// i2i / inpaint 互斥
	if req.I2I != nil && req.Inpaint != nil {
		return newValidationError("", "i2i 与 inpaint 互斥，不能同时使用")
	}

	// controlnet / character_references 互斥
	if req.Controlnet != nil && len(req.CharacterRefs) > 0 {
		return newValidationError("", "controlnet 与 character_references 互斥，不能同时使用")
	}

	// i2i 校验
	if req.I2I != nil {
		if err := validateImageRefField(req.I2I.Image, "i2i.image"); err != nil {
			return err
		}
		if err := validateStrength(req.I2I.Strength, 0.01, 0.99, "i2i.strength"); err != nil {
			return err
		}
		if req.I2I.Noise != nil {
			if err := validateRange(req.I2I.Noise, 0.0, 0.99, "i2i.noise"); err != nil {
				return err
			}
		}
		// 尺寸一致性
		if sizeW > 0 {
			if err := validateImageSizeMatch(req.I2I.Image, sizeW, sizeH, "i2i.image"); err != nil {
				return err
			}
		}
	}

	// inpaint 校验
	if req.Inpaint != nil {
		if err := validateImageRefField(req.Inpaint.Image, "inpaint.image"); err != nil {
			return err
		}
		if err := validateImageRefField(req.Inpaint.Mask, "inpaint.mask"); err != nil {
			return err
		}
		if err := validateStrength(req.Inpaint.Strength, 0.01, 1.0, "inpaint.strength"); err != nil {
			return err
		}
		if sizeW > 0 {
			if err := validateImageSizeMatch(req.Inpaint.Image, sizeW, sizeH, "inpaint.image"); err != nil {
				return err
			}
			if err := validateImageSizeMatch(req.Inpaint.Mask, sizeW, sizeH, "inpaint.mask"); err != nil {
				return err
			}
		}
	}

	// controlnet 校验
	if req.Controlnet != nil {
		if len(req.Controlnet.Images) == 0 {
			return newValidationError("controlnet.images", "至少需要 1 张参考图")
		}
		if len(req.Controlnet.Images) > 4 {
			return newValidationError("controlnet.images", "最多 4 张参考图")
		}
		for i, img := range req.Controlnet.Images {
			prefix := "controlnet.images[" + strconv.Itoa(i) + "]"
			// cache_id 复用态 或 image 完整态，二选一
			if img.CacheID != "" {
				if img.Image != "" || img.InfoExtracted != nil {
					return newValidationError(prefix, "cache_id 复用态不允许混合 image 或 info_extracted")
				}
			} else if img.Image != "" {
				if err := validateImageRefField(img.Image, prefix+".image"); err != nil {
					return err
				}
				if img.InfoExtracted != nil {
					if err := validateRange(img.InfoExtracted, 0.01, 1.0, prefix+".info_extracted"); err != nil {
						return err
					}
				}
			} else {
				return newValidationError(prefix, "必须提供 image 或 cache_id")
			}
			if img.Strength != nil {
				if err := validateRange(img.Strength, 0.01, 1.0, prefix+".strength"); err != nil {
					return err
				}
			}
		}
		if req.Controlnet.Strength != nil {
			if err := validateRange(req.Controlnet.Strength, 0.0, 1.0, "controlnet.strength"); err != nil {
				return err
			}
		}
	}

	// character_references 校验
	if len(req.CharacterRefs) > 0 {
		if len(req.CharacterRefs) > 1 {
			return newValidationError("character_references", "最多 1 张角色参考图")
		}
		for i, cr := range req.CharacterRefs {
			prefix := "character_references[" + strconv.Itoa(i) + "]"
			if err := validateImageRefField(cr.Image, prefix+".image"); err != nil {
				return err
			}
			if cr.Type != "" && !allowedCharRefTypes[cr.Type] {
				return newValidationError(prefix+".type", "type 只能是 character / style / character&style")
			}
			if cr.Fidelity != nil {
				if err := validateRange(cr.Fidelity, 0.0, 1.0, prefix+".fidelity"); err != nil {
					return err
				}
			}
			if cr.Strength != nil {
				if err := validateRange(cr.Strength, 0.0, 1.0, prefix+".strength"); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func checkNoCJK(field, s string) error {
	if s == "" {
		return nil
	}
	if cjkRe.MatchString(s) {
		return newValidationError(field, "不允许包含中文 / 日文假名 / 韩文 / 全角符号，请使用英文")
	}
	return nil
}

func validateSize(w, h int) error {
	if w < 256 || h < 256 {
		return newValidationError("size", "宽高不能小于 256")
	}
	if w > maxWidth || h > maxHeight {
		return newValidationError("size", "宽高不能超过 1216")
	}
	if w%64 != 0 || h%64 != 0 {
		return newValidationError("size", "宽高必须是 64 的倍数")
	}
	if w*h > maxPixels {
		return newValidationError("size", "总像素不能超过 1024x1024")
	}
	return nil
}

func validateImageRefField(s, field string) error {
	if strings.TrimSpace(s) == "" {
		return newValidationError(field, "图片不能为空")
	}
	return nil
}

func validateStrength(v *float64, min, max float64, field string) error {
	if v == nil {
		return nil
	}
	return validateRange(v, min, max, field)
}

func validateRange(v *float64, min, max float64, field string) error {
	if v == nil {
		return nil
	}
	if *v < min || *v > max {
		return newValidationError(field, strconv.FormatFloat(*v, 'f', -1, 64)+" 超出范围 "+strconv.FormatFloat(min, 'f', -1, 64)+"~"+strconv.FormatFloat(max, 'f', -1, 64))
	}
	return nil
}

// validateImageSizeMatch 解码图片并校验宽高与 size 一致。解码失败时跳过（交给上游校验）。
func validateImageSizeMatch(dataURI string, expectW, expectH int, field string) error {
	_, data, err := DecodeDataURI(dataURI)
	if err != nil {
		return newValidationError(field, "图片解码失败: "+err.Error())
	}
	w, h, err := DecodeImageMeta(data)
	if err != nil {
		// 不支持的格式（如 webp），跳过尺寸校验
		return nil
	}
	if w != expectW || h != expectH {
		return newValidationError(field, "图片宽高 "+strconv.Itoa(w)+"x"+strconv.Itoa(h)+" 与 size "+strconv.Itoa(expectW)+"x"+strconv.Itoa(expectH)+" 不一致")
	}
	return nil
}
