export const DEFAULT_VIBE_INFO_EXTRACTED = 0.7
export const DEFAULT_VIBE_STRENGTH = 0.6
export const DEFAULT_CONTROLNET_STRENGTH = 1

export function clampVibeInfoExtracted(v: number) {
  if (!Number.isFinite(v)) return DEFAULT_VIBE_INFO_EXTRACTED
  return Math.min(1, Math.max(0.01, Math.round(v * 100) / 100))
}

export function clampVibeStrength(v: number) {
  if (!Number.isFinite(v)) return DEFAULT_VIBE_STRENGTH
  return Math.min(1, Math.max(0.01, Math.round(v * 100) / 100))
}

export function clampControlnetStrength(v: number) {
  if (!Number.isFinite(v)) return DEFAULT_CONTROLNET_STRENGTH
  return Math.min(1, Math.max(0, Math.round(v * 100) / 100))
}

/** 图库：1 张改图 + 最多 4 张 Vibe 参考 */
export const NAI_GALLERY_MAX_INPUT_IMAGES = 5
