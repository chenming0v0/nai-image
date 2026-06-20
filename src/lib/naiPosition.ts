import type { NaiCharacterPosition } from '../types'

const POSITION_RE = /^[A-E][1-5]$/

export function isNaiPosition(value: string): value is NaiCharacterPosition {
  return POSITION_RE.test(value)
}

export function normalizeNaiPosition(value: string | undefined | null): NaiCharacterPosition {
  if (value && isNaiPosition(value)) return value
  return 'C3'
}
