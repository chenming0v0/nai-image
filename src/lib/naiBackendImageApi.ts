import type { ApiProfile } from '../types'
import type { CallApiOptions, CallApiResult } from './imageApiShared'
import { normalizeBase64Image } from './imageApiShared'
import { normalizeBaseUrl } from './devProxy'

function parseSize(size: string): [number, number] | undefined {
  const match = size.match(/^(\d+)x(\d+)$/)
  if (!match) return undefined
  return [Number(match[1]), Number(match[2])]
}

function numberOrUndefined(value: unknown): number | undefined {
  return typeof value === 'number' && Number.isFinite(value) ? value : undefined
}

function nullableNumberOrUndefined(value: unknown): number | undefined {
  return value == null ? undefined : numberOrUndefined(value)
}

function boolOrUndefined(value: unknown): boolean | undefined {
  return typeof value === 'boolean' ? value : undefined
}

function dataUrlToBackendUrl(url: string) {
  return url
}

async function syncProfileToBackend(profile: ApiProfile) {
  const baseUrl = profile.baseUrl.trim()
  const apiKey = profile.apiKey.trim()
  if (!baseUrl) throw new Error('请先填写 API 地址')

  const response = await fetch('http://127.0.0.1:8787/api/settings', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      upstream_base_url: normalizeBaseUrl(baseUrl),
      upstream_api_key: apiKey,
      default_model: profile.model,
      request_timeout_seconds: profile.timeout,
    }),
  })
  if (!response.ok) throw new Error(`保存后端配置失败 HTTP ${response.status}`)
}

export async function callNaiBackendImageApi(opts: CallApiOptions, profile: ApiProfile): Promise<CallApiResult> {
  await syncProfileToBackend(profile)

  const params = opts.params
  const inputImages = opts.inputImageDataUrls
  const referenceImages = inputImages.slice(opts.maskDataUrl ? 1 : inputImages.length > 0 ? 1 : 0, 5)
  const body: Record<string, unknown> = {
    model: profile.model,
    prompt: opts.prompt,
    size: parseSize(params.size),
    steps: numberOrUndefined(params.steps),
    scale: numberOrUndefined(params.scale),
    sampler: params.sampler,
    seed: nullableNumberOrUndefined(params.seed),
    variety_boost: boolOrUndefined(params.variety_boost),
    cfg_rescale: nullableNumberOrUndefined(params.cfg_rescale),
    noise_schedule: params.noise_schedule,
    image_format: params.image_format ?? params.output_format,
    n_samples: 1,
    max_tokens: 4096,
    stream: false,
  }

  const extra = opts as CallApiOptions & {
    negativePrompt?: string
    characters?: Array<{ prompt: string; negative_prompt?: string; position?: string }>
    use_coords?: boolean
    use_order?: boolean
  }
  if (extra.negativePrompt?.trim()) body.negative_prompt = extra.negativePrompt.trim()
  if (extra.characters?.length) {
    body.characters = extra.characters
      .filter((character) => character.prompt.trim())
      .map((character) => ({
        prompt: character.prompt.trim(),
        negative_prompt: character.negative_prompt?.trim() || undefined,
        position: character.position,
      }))
  }
  if (extra.use_coords !== undefined) body.use_coords = extra.use_coords
  if (extra.use_order !== undefined) body.use_order = extra.use_order

  if (opts.maskDataUrl && inputImages[0]) {
    body.inpaint = {
      image: dataUrlToBackendUrl(inputImages[0]),
      mask: dataUrlToBackendUrl(opts.maskDataUrl),
      strength: nullableNumberOrUndefined(params.inpaint_strength),
    }
  } else if (inputImages[0]) {
    body.i2i = {
      image: dataUrlToBackendUrl(inputImages[0]),
      strength: nullableNumberOrUndefined(params.i2i_strength),
      noise: nullableNumberOrUndefined(params.i2i_noise),
    }
  }

  if (referenceImages.length) {
    body.controlnet = {
      strength: nullableNumberOrUndefined(params.controlnet_strength),
      images: referenceImages.map((image) => ({ image: dataUrlToBackendUrl(image) })),
    }
  }

  const response = await fetch('http://127.0.0.1:8787/api/generate', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(body),
  })
  const payload = await response.json().catch(() => null)
  if (!response.ok) {
    const message = payload && typeof payload === 'object' && 'error' in payload ? String(payload.error) : `HTTP ${response.status}`
    throw new Error(message)
  }

  const imageUrl = payload && typeof payload === 'object' && 'image' in payload && payload.image && typeof payload.image === 'object' && 'url' in payload.image
    ? String(payload.image.url)
    : ''
  if (!imageUrl) throw new Error('后端没有返回图片')

  const imageResponse = await fetch(`http://127.0.0.1:8787${imageUrl}`)
  if (!imageResponse.ok) throw new Error(`读取后端图片失败 HTTP ${imageResponse.status}`)
  const blob = await imageResponse.blob()
  const dataUrl = await new Promise<string>((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(String(reader.result))
    reader.onerror = () => reject(reader.error ?? new Error('读取图片失败'))
    reader.readAsDataURL(blob)
  })

  return {
    images: [normalizeBase64Image(dataUrl, blob.type || 'image/png')],
    actualParams: {
      ...params,
      n: 1,
    },
  }
}
