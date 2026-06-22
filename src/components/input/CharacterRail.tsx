import { AnimatePresence, motion } from 'framer-motion'
import type { CharacterSlot } from '../../types'
import { characterColor } from '../../lib/characterColors'
import PositionGrid from './PositionGrid'

export type CharacterGridItem = {
  id: string
  position: CharacterSlot['position']
  color: string
  label: string
}

type Props = {
  characters: CharacterSlot[]
  useCoords: boolean
  useOrder: boolean
  onAddCharacter: () => void
  onUpdateCharacter: (id: string, patch: Partial<CharacterSlot>) => void
  onRemoveCharacter: (id: string) => void
  onUseCoordsChange: (v: boolean) => void
  onUseOrderChange: (v: boolean) => void
}

export function getCharacterGridItems(characters: CharacterSlot[]): CharacterGridItem[] {
  return characters.map((c, i) => ({
    id: c.id,
    position: c.position,
    color: characterColor(i),
    label: c.prompt.slice(0, 24) || `Character ${i + 1}`,
  }))
}

export default function CharacterRail({
  characters,
  useCoords,
  useOrder,
  onAddCharacter,
  onUpdateCharacter,
  onRemoveCharacter,
  onUseCoordsChange,
  onUseOrderChange,
}: Props) {
  return (
    <div className="w-full h-full flex flex-col">
      {/* 标题栏 - 固定 */}
      <div className="flex-shrink-0 flex items-center justify-between gap-2 mb-3">
        <span className="text-xs font-semibold tracking-wide text-gray-700 dark:text-gray-200">Characters</span>
        <button
          type="button"
          onClick={onAddCharacter}
          className="rounded-lg border border-gray-300/70 bg-white/55 px-2 py-1 text-[11px] text-gray-700 shadow-sm transition hover:bg-white/80 dark:border-white/[0.10] dark:bg-white/[0.06] dark:text-gray-200 dark:hover:bg-white/[0.10]"
        >
          + 添加角色
        </button>
      </div>

      {/* 选项栏 - 固定 */}
      <div className="flex-shrink-0 flex flex-wrap items-center gap-2 mb-3">
        <label className="group flex cursor-pointer items-center gap-2 rounded-xl border border-gray-300/40 bg-white/40 px-3 py-1.5 text-[11px] text-gray-700 shadow-sm backdrop-blur-md transition-all hover:border-cyan-400/50 hover:bg-white/50 dark:border-white/[0.12] dark:bg-white/[0.06] dark:text-gray-300 dark:hover:border-cyan-400/40 dark:hover:bg-white/[0.10] has-[:checked]:border-cyan-400/60 has-[:checked]:bg-cyan-50/50 has-[:checked]:text-cyan-700 has-[:checked]:backdrop-blur-xl dark:has-[:checked]:border-cyan-400/50 dark:has-[:checked]:bg-cyan-500/15 dark:has-[:checked]:text-cyan-300">
          <input
            type="checkbox"
            checked={useCoords}
            onChange={(e) => onUseCoordsChange(e.target.checked)}
            className="h-3.5 w-3.5 rounded border-gray-300 text-cyan-600 transition focus:ring-2 focus:ring-cyan-500/20 dark:border-white/[0.15] dark:bg-white/[0.05] dark:checked:bg-cyan-500 dark:checked:border-cyan-500"
          />
          使用坐标
        </label>
        <label className="group flex cursor-pointer items-center gap-2 rounded-xl border border-gray-300/40 bg-white/40 px-3 py-1.5 text-[11px] text-gray-700 shadow-sm backdrop-blur-md transition-all hover:border-cyan-400/50 hover:bg-white/50 dark:border-white/[0.12] dark:bg-white/[0.06] dark:text-gray-300 dark:hover:border-cyan-400/40 dark:hover:bg-white/[0.10] has-[:checked]:border-cyan-400/60 has-[:checked]:bg-cyan-50/50 has-[:checked]:text-cyan-700 has-[:checked]:backdrop-blur-xl dark:has-[:checked]:border-cyan-400/50 dark:has-[:checked]:bg-cyan-500/15 dark:has-[:checked]:text-cyan-300">
          <input
            type="checkbox"
            checked={useOrder}
            onChange={(e) => onUseOrderChange(e.target.checked)}
            className="h-3.5 w-3.5 rounded border-gray-300 text-cyan-600 transition focus:ring-2 focus:ring-cyan-500/20 dark:border-white/[0.15] dark:bg-white/[0.05] dark:checked:bg-cyan-500 dark:checked:border-cyan-500"
          />
          按顺序排序
        </label>
      </div>

      {/* 角色列表 - 可滚动 */}
      <div className="flex-1 min-h-0 overflow-y-auto hide-scrollbar">
        {characters.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-gray-300/80 bg-white/35 px-3 py-6 text-center text-[11px] text-gray-500 shadow-inner dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-500">
            <p>暂无角色</p>
            <p className="mt-1">点击添加，用于多角色构图与坐标</p>
          </div>
        ) : (
          <div className="space-y-2 pr-0.5">
            <AnimatePresence initial={false}>
              {characters.map((c, i) => (
                <motion.div
                  key={c.id}
                  layout
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.90, transition: { duration: 0.2 } }}
                  transition={{ type: 'spring', stiffness: 350, damping: 28 }}
                  className="relative overflow-hidden rounded-2xl border border-gray-300/70 bg-white/58 p-3 shadow-[0_10px_28px_rgba(15,23,42,0.08)] ring-1 ring-white/50 backdrop-blur-xl dark:border-white/[0.10] dark:bg-white/[0.05] dark:ring-white/[0.04]"
                >
                  <div className="absolute bottom-0 left-0 top-0 w-1 rounded-l-xl" style={{ background: characterColor(i) }} />
                  <div className="mb-2 flex items-center justify-between pl-2">
                    <div className="flex items-center gap-1.5">
                      <span className="h-2 w-2 rounded-full" style={{ background: characterColor(i) }} />
                      <span className="font-mono text-[10px] uppercase tracking-wider text-gray-600 dark:text-gray-300">Character {i + 1}</span>
                    </div>
                    <button type="button" onClick={() => onRemoveCharacter(c.id)} className="text-gray-400 hover:text-red-500 transition" title="删除">
                      ×
                    </button>
                  </div>
                  <textarea
                    value={c.prompt}
                    onChange={(e) => onUpdateCharacter(c.id, { prompt: e.target.value })}
                    placeholder="1girl, blue hair, blue dress…"
                    rows={2}
                    className="mb-1.5 w-full resize-y rounded-xl border border-gray-300/60 bg-white/70 px-2.5 py-2 text-xs shadow-inner outline-none transition focus:border-violet-300 focus:bg-white/85 dark:border-white/[0.08] dark:bg-white/[0.04]"
                  />
                  <textarea
                    value={c.negative_prompt}
                    onChange={(e) => onUpdateCharacter(c.id, { negative_prompt: e.target.value })}
                    placeholder="negative（可选）"
                    rows={1}
                    className="w-full resize-y rounded-xl border border-gray-300/60 bg-white/70 px-2.5 py-2 text-xs shadow-inner outline-none transition focus:border-violet-300 focus:bg-white/85 dark:border-white/[0.08] dark:bg-white/[0.04]"
                  />
                  {useCoords && (
                    <div className="mt-2 flex justify-end">
                      <PositionGrid compact value={c.position} onChange={(p) => onUpdateCharacter(c.id, { position: p })} />
                    </div>
                  )}
                </motion.div>
              ))}
            </AnimatePresence>
          </div>
        )}
      </div>
    </div>
  )
}
