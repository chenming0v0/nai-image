import { DEFAULT_PARAMS, type AppSettings, type TaskParams } from '../types'
import { getActiveApiProfile } from './apiProfiles'
import { normalizeImageSize } from './size'
import { normalizeNaiSize } from './naiSizes'
import { clampNaiScale, clampNaiSteps } from './naiDrawParams'

export const DEFAULT_FAL_IMAGE_SIZE = '1360x1024'
export const MAX_FAL_OUTPUT_IMAGES = 4
export const MAX_OPENAI_OUTPUT_IMAGES = 10

export function getOutputImageLimitForSettings(settings: AppSettings) {
  return getActiveApiProfile(settings).provider === 'fal' ? MAX_FAL_OUTPUT_IMAGES : MAX_OPENAI_OUTPUT_IMAGES
}

export function normalizeParamsForSettings(
  params: TaskParams,
  settings: AppSettings,
  options: { hasInputImages?: boolean } = {},
): TaskParams {
  const activeProfile = getActiveApiProfile(settings)
  const outputImageLimit = getOutputImageLimitForSettings(settings)
  const nextParams: TaskParams = {
    ...params,
    size: activeProfile.provider === 'fal'
      ? (normalizeImageSize(params.size) || DEFAULT_PARAMS.size)
      : normalizeNaiSize(params.size),
    n: Math.min(outputImageLimit, Math.max(1, params.n || DEFAULT_PARAMS.n)),
  }

  if (activeProfile.provider === 'openai' && activeProfile.codexCli) {
    nextParams.quality = DEFAULT_PARAMS.quality
  }

  if (activeProfile.provider === 'fal') {
    if (!options.hasInputImages && nextParams.size === 'auto') nextParams.size = DEFAULT_FAL_IMAGE_SIZE
    if (nextParams.quality === 'auto') nextParams.quality = 'high'
    nextParams.moderation = DEFAULT_PARAMS.moderation
    nextParams.output_compression = DEFAULT_PARAMS.output_compression
  }


  if (activeProfile.provider !== 'fal') {
    nextParams.size = normalizeNaiSize(nextParams.size)
    nextParams.steps = clampNaiSteps(nextParams.steps ?? DEFAULT_PARAMS.steps ?? 23)
    nextParams.scale = clampNaiScale(nextParams.scale ?? DEFAULT_PARAMS.scale ?? 5)
    if (!nextParams.sampler) nextParams.sampler = DEFAULT_PARAMS.sampler
    if (!nextParams.noise_schedule) nextParams.noise_schedule = DEFAULT_PARAMS.noise_schedule
    if (!nextParams.image_format) nextParams.image_format = DEFAULT_PARAMS.image_format
    if (nextParams.seed === undefined) nextParams.seed = null
    if (nextParams.variety_boost === undefined) nextParams.variety_boost = false
  }

  if (nextParams.output_format === 'png') {
    nextParams.output_compression = DEFAULT_PARAMS.output_compression
  }

  return nextParams
}

export function getChangedParams(current: TaskParams, next: TaskParams): Partial<TaskParams> {
  const patch: Partial<TaskParams> = {}
  for (const key of Object.keys(next) as Array<keyof TaskParams>) {
    if (current[key] !== next[key]) {
      ;(patch as Record<keyof TaskParams, TaskParams[keyof TaskParams]>)[key] = next[key]
    }
  }
  return patch
}
