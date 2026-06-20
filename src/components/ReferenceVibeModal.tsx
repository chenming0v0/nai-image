import { useEffect, useRef, useState } from 'react'
import type { InputImage } from '../types'
import { usePreventBackgroundScroll } from '../hooks/usePreventBackgroundScroll'
import {
  DEFAULT_VIBE_INFO_EXTRACTED,
  DEFAULT_VIBE_STRENGTH,
  clampVibeInfoExtracted,
  clampVibeStrength,
} from '../lib/naiVibeDefaults'

type Props = {
  image: InputImage
  referenceIndex: number
  onClose: () => void
  onSave: (patch: Pick<InputImage, 'vibeInfoExtracted' | 'vibeStrength'>) => void
  onReplace: () => void
}

export default function ReferenceVibeModal({ image, referenceIndex, onClose, onSave, onReplace }: Props) {
  usePreventBackgroundScroll(true)
  const modalRef = useRef<HTMLDivElement>(null)
  const mouseDownTargetRef = useRef<EventTarget | null>(null)

  const [info, setInfo] = useState(image.vibeInfoExtracted ?? DEFAULT_VIBE_INFO_EXTRACTED)
  const [strength, setStrength] = useState(image.vibeStrength ?? DEFAULT_VIBE_STRENGTH)

  useEffect(() => {
    setInfo(image.vibeInfoExtracted ?? DEFAULT_VIBE_INFO_EXTRACTED)
    setStrength(image.vibeStrength ?? DEFAULT_VIBE_STRENGTH)
  }, [image.id, image.vibeInfoExtracted, image.vibeStrength])

  const handleBackdropMouseUp = (e: React.MouseEvent) => {
    const down = mouseDownTargetRef.current
    if (modalRef.current && down && !modalRef.current.contains(down as Node) && !modalRef.current.contains(e.target as Node)) {
      onClose()
    }
    mouseDownTargetRef.current = null
  }

  const save = () => {
    onSave({
      vibeInfoExtracted: clampVibeInfoExtracted(info),
      vibeStrength: clampVibeStrength(strength),
    })
    onClose()
  }

  return (
    <div
      data-no-drag-select
      className="fixed inset-0 z-[75] flex items-center justify-center p-4"
      onMouseDown={(e) => { mouseDownTargetRef.current = e.target }}
      onMouseUp={handleBackdropMouseUp}
    >
      <div className="absolute inset-0 bg-black/30 backdrop-blur-sm animate-overlay-in" />
      <div
        ref={modalRef}
        className="relative z-10 w-full max-w-md rounded-3xl border border-white/50 bg-white/95 p-5 shadow-2xl ring-1 ring-black/5 animate-modal-in dark:border-white/[0.08] dark:bg-gray-900/95 dark:ring-white/10"
      >
        <div className="mb-4 flex items-start justify-between gap-3">
          <div>
            <h3 className="text-base font-semibold text-gray-800 dark:text-gray-100">Vibe 参考 #{referenceIndex}</h3>
            <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">对应 controlnet.images 的 info_extracted 与 strength</p>
          </div>
          <button type="button" onClick={onClose} className="rounded-full p-1 text-gray-400 hover:bg-gray-100 dark:hover:bg-white/[0.06]" aria-label="关闭">
            <svg className="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
          </button>
        </div>

        <div className="mb-4 flex gap-3">
          <div className="h-20 w-20 shrink-0 overflow-hidden rounded-xl border border-gray-200/70 dark:border-white/[0.08]">
            {image.dataUrl && <img src={image.dataUrl} alt="" className="h-full w-full object-cover" />}
          </div>
          <p className="text-xs leading-relaxed text-gray-500 dark:text-gray-400">
            参考图尺寸无需与输出 size 一致。上传后由上游做 resize。若已有 cache_id 可在后续版本在此填写以省流量费。
          </p>
        </div>

        <label className="mb-4 block">
          <div className="mb-1 flex justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>信息提取 info_extracted</span>
            <span className="font-mono">{clampVibeInfoExtracted(info).toFixed(2)}</span>
          </div>
          <input type="range" min={0.01} max={1} step={0.01} value={info} onChange={(e) => setInfo(Number(e.target.value))} className="w-full accent-blue-500" />
        </label>

        <label className="mb-5 block">
          <div className="mb-1 flex justify-between text-xs text-gray-500 dark:text-gray-400">
            <span>参考强度 strength</span>
            <span className="font-mono">{clampVibeStrength(strength).toFixed(2)}</span>
          </div>
          <input type="range" min={0.01} max={1} step={0.01} value={strength} onChange={(e) => setStrength(Number(e.target.value))} className="w-full accent-violet-500" />
        </label>

        <div className="flex flex-wrap gap-2">
          <button type="button" onClick={onReplace} className="rounded-xl border border-gray-200/70 px-4 py-2 text-sm text-gray-600 hover:bg-gray-50 dark:border-white/[0.08] dark:text-gray-300 dark:hover:bg-white/[0.06]">
            替换图片
          </button>
          <button type="button" onClick={onClose} className="rounded-xl border border-gray-200/70 px-4 py-2 text-sm text-gray-600 dark:border-white/[0.08] dark:text-gray-300">
            取消
          </button>
          <button type="button" onClick={save} className="ml-auto rounded-xl bg-blue-500 px-4 py-2 text-sm font-medium text-white hover:bg-blue-600">
            保存
          </button>
        </div>
      </div>
    </div>
  )
}
