import { motion } from 'framer-motion'
import type { NaiCharacterPosition } from '../../types'
import { normalizeNaiPosition } from '../../lib/naiPosition'

const COLS = ['A', 'B', 'C', 'D', 'E'] as const
const ROWS = [1, 2, 3, 4, 5] as const

export type PositionGridCharacter = {
  id: string
  position: NaiCharacterPosition
  color: string
  label?: string
}

type Props = {
  value: NaiCharacterPosition
  onChange: (p: NaiCharacterPosition) => void
  characters?: PositionGridCharacter[]
  compact?: boolean
  className?: string
}

export default function PositionGrid({ value, onChange, characters = [], compact, className = '' }: Props) {
  const active = normalizeNaiPosition(value)
  return (
    <div className={className}>
      <div
        className={`relative grid grid-cols-5 gap-1 rounded-2xl border border-gray-300/70 bg-white/52 p-2 shadow-inner ring-1 ring-white/45 dark:border-white/[0.08] dark:bg-white/[0.04] aspect-square ${
          compact ? 'max-w-[168px]' : 'max-w-[220px]'
        }`}
      >
        {ROWS.map((r) =>
          COLS.map((c) => {
            const pos = `${c}${r}` as NaiCharacterPosition
            const isActive = active === pos
            const charsHere = characters.filter((x) => normalizeNaiPosition(x.position) === pos)
            return (
              <button
                key={pos}
                type="button"
                onClick={() => onChange(pos)}
                aria-label={`方位 ${pos}`}
                className={`relative flex aspect-square items-center justify-center rounded-md border text-[10px] font-mono transition-colors ${
                  isActive
                    ? 'border-violet-500/80 bg-violet-500/15 text-violet-600 dark:text-violet-300'
                    : 'border-gray-300/60 bg-white/35 text-gray-500 hover:border-gray-400 hover:bg-white/60 dark:border-white/[0.07] dark:bg-white/[0.02] dark:hover:bg-white/[0.05]'
                }`}
              >
                <span className="opacity-60">{pos}</span>
                {charsHere.length > 0 && (
                  <div className="pointer-events-none absolute inset-0 flex items-center justify-center gap-0.5">
                    {charsHere.map((ch) => (
                      <motion.span
                        key={ch.id}
                        layoutId={`nai-char-dot-${ch.id}`}
                        className="rounded-full shadow-sm"
                        style={{
                          backgroundColor: ch.color,
                          width: compact ? 8 : 10,
                          height: compact ? 8 : 10,
                          boxShadow: `0 0 8px ${ch.color}`,
                        }}
                        title={ch.label}
                        transition={{ type: 'spring', stiffness: 500, damping: 30 }}
                      />
                    ))}
                  </div>
                )}
                {isActive && (
                  <motion.div
                    layoutId="nai-active-pos-ring"
                    className="pointer-events-none absolute inset-0 rounded-md ring-2 ring-violet-500/50"
                    transition={{ type: 'spring', stiffness: 500, damping: 35 }}
                  />
                )}
              </button>
            )
          }),
        )}
      </div>
      {!compact && (
        <div className="mt-1.5 flex justify-between font-mono text-[10px] text-gray-400 dark:text-gray-500">
          <span>左</span>
          <span>{active}</span>
          <span>右</span>
        </div>
      )}
    </div>
  )
}
