export const CHARACTER_COLORS = [
  '#8b5cf6',
  '#06b6d4',
  '#f59e0b',
  '#ec4899',
  '#22c55e',
  '#3b82f6',
  '#ef4444',
  '#a855f7',
] as const

export function characterColor(index: number) {
  return CHARACTER_COLORS[index % CHARACTER_COLORS.length]
}
