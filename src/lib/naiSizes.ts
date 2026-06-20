export type NaiSizePresetId = 'portrait' | 'landscape' | 'square'

export const NAI_SIZE_PRESETS: Array<{
  id: NaiSizePresetId
  label: string
  size: string
  width: number
  height: number
}> = [
  { id: 'portrait', label: '竖图', size: '832x1216', width: 832, height: 1216 },
  { id: 'landscape', label: '横图', size: '1216x832', width: 1216, height: 832 },
  { id: 'square', label: '方图', size: '1024x1024', width: 1024, height: 1024 },
]

export const DEFAULT_NAI_SIZE = NAI_SIZE_PRESETS[0].size

export function normalizeNaiSize(size: string | undefined | null): string {
  const raw = (size || '').trim().toLowerCase().replace(/×/g, 'x')
  for (const p of NAI_SIZE_PRESETS) {
    if (raw === p.size) return p.size
  }
  const m = raw.match(/^(\d+)\s*x\s*(\d+)$/)
  if (m) {
    const w = Number(m[1])
    const h = Number(m[2])
    for (const p of NAI_SIZE_PRESETS) {
      if (p.width === w && p.height === h) return p.size
    }
  }
  return DEFAULT_NAI_SIZE
}

export function formatNaiSizeLabel(size: string): string {
  const normalized = normalizeNaiSize(size)
  const preset = NAI_SIZE_PRESETS.find((p) => p.size === normalized)
  if (!preset) return normalized
  return `${preset.label} ${preset.width}×${preset.height}`
}
