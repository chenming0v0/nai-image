import { useRef } from 'react'
import { usePreventBackgroundScroll } from '../hooks/usePreventBackgroundScroll'
import { NAI_SIZE_PRESETS, normalizeNaiSize } from '../lib/naiSizes'

interface Props {
  currentSize: string
  onSelect: (size: string) => void
  onClose: () => void
  allowAuto?: boolean
}

export default function SizePickerModal({ currentSize, onSelect, onClose }: Props) {
  usePreventBackgroundScroll(true)

  const modalRef = useRef<HTMLDivElement>(null)
  const mouseDownTargetRef = useRef<EventTarget | null>(null)
  const active = normalizeNaiSize(currentSize)

  const handleMouseDown = (e: React.MouseEvent) => {
    mouseDownTargetRef.current = e.target
  }

  const handleMouseUp = (e: React.MouseEvent) => {
    const mouseDownTarget = mouseDownTargetRef.current
    const mouseUpTarget = e.target
    if (
      modalRef.current &&
      mouseDownTarget &&
      !modalRef.current.contains(mouseDownTarget as Node) &&
      mouseUpTarget &&
      !modalRef.current.contains(mouseUpTarget as Node)
    ) {
      onClose()
    }
    mouseDownTargetRef.current = null
  }

  return (
    <div
      data-no-drag-select
      className="fixed inset-0 z-[70] flex items-center justify-center p-4"
      onMouseDown={handleMouseDown}
      onMouseUp={handleMouseUp}
    >
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm animate-overlay-in" />
      <div
        ref={modalRef}
        className="relative z-10 w-full max-w-sm rounded-3xl border border-white/50 bg-white/95 p-5 shadow-2xl ring-1 ring-black/5 animate-modal-in dark:border-white/[0.08] dark:bg-gray-900/95 dark:ring-white/10"
      >
        <div className="mb-4 flex items-start justify-between gap-4">
          <div>
            <h3 className="text-base font-semibold text-gray-800 dark:text-gray-100">图像尺寸</h3>
            <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">仅支持文档规定的三种分辨率</p>
          </div>
          <button
            type="button"
            onClick={onClose}
            className="rounded-full p-1 text-gray-400 transition hover:bg-gray-100 hover:text-gray-600 dark:hover:bg-white/[0.06] dark:hover:text-gray-200"
            aria-label="关闭"
          >
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div className="grid gap-2">
          {NAI_SIZE_PRESETS.map((preset) => {
            const selected = active === preset.size
            const isLandscape = preset.width > preset.height
            const isSquare = preset.width === preset.height
            return (
              <button
                key={preset.id}
                type="button"
                onClick={() => {
                  onSelect(preset.size)
                  onClose()
                }}
                className={`flex items-center gap-3 rounded-2xl border px-4 py-3 text-left transition ${
                  selected
                    ? 'border-blue-400 bg-blue-50 text-blue-700 dark:border-blue-500/50 dark:bg-blue-500/10 dark:text-blue-200'
                    : 'border-gray-200/70 bg-white/60 text-gray-700 hover:bg-gray-50 dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-200 dark:hover:bg-white/[0.06]'
                }`}
              >
                <div
                  className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border-2 ${
                    selected ? 'border-current' : 'border-gray-300 dark:border-white/20'
                  }`}
                >
                  <div
                    className="rounded-[2px] bg-current opacity-50"
                    style={{
                      width: isSquare ? 18 : isLandscape ? 22 : 14,
                      height: isSquare ? 18 : isLandscape ? 14 : 22,
                    }}
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="text-sm font-medium">{preset.label}</div>
                  <div className="font-mono text-xs opacity-70">
                    {preset.width} × {preset.height}
                  </div>
                </div>
              </button>
            )
          })}
        </div>
      </div>
    </div>
  )
}
