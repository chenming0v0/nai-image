export const NAI_SAMPLER_OPTIONS = [
  { label: 'Euler Ancestral', value: 'k_euler_ancestral' },
  { label: 'Euler', value: 'k_euler' },
  { label: 'DPM++ 2M', value: 'k_dpmpp_2m' },
  { label: 'DPM++ SDE', value: 'k_dpmpp_sde' },
  { label: 'DPM++ 2S Ancestral', value: 'k_dpmpp_2s_ancestral' },
  { label: 'DPM 2', value: 'k_dpm_2' },
  { label: 'DPM 2 Ancestral', value: 'k_dpm_2_ancestral' },
  { label: 'DDIM', value: 'ddim' },
] as const

export const NAI_NOISE_SCHEDULE_OPTIONS = [
  { label: 'Karras', value: 'karras' },
  { label: 'Exponential', value: 'exponential' },
  { label: 'Polyexponential', value: 'polyexponential' },
] as const

export const NAI_STEPS_MIN = 1
export const NAI_STEPS_MAX = 28
export const NAI_STEPS_DEFAULT = 23
export const NAI_SCALE_DEFAULT = 5

export function clampNaiSteps(value: number) {
  if (!Number.isFinite(value)) return NAI_STEPS_DEFAULT
  return Math.min(NAI_STEPS_MAX, Math.max(NAI_STEPS_MIN, Math.round(value)))
}

export function clampNaiScale(value: number) {
  if (!Number.isFinite(value)) return NAI_SCALE_DEFAULT
  return Math.min(20, Math.max(1, Math.round(value * 10) / 10))
}

export function parseNaiSeedInput(raw: string): number | null {
  const text = raw.trim()
  if (!text) return null
  const n = Number(text)
  if (!Number.isFinite(n) || !Number.isInteger(n) || n < 0) return null
  return n
}
