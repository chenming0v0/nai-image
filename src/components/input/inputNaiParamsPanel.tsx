import { useState } from 'react'
import type { TaskParams } from '../../types'
import { dismissAllTooltips } from '../../lib/tooltipDismiss'
import Select from '../Select'
import {
  NAI_NOISE_SCHEDULE_OPTIONS,
  NAI_SAMPLER_OPTIONS,
  NAI_SCALE_DEFAULT,
  NAI_STEPS_DEFAULT,
  NAI_STEPS_MAX,
  NAI_STEPS_MIN,
  clampNaiScale,
  clampNaiSteps,
  parseNaiSeedInput,
} from '../../lib/naiDrawParams'

type Props = {
  cols: string
  params: TaskParams
  setParams: (patch: Partial<TaskParams>) => void
  displaySizeLabel?: string
  displaySize: string
  selectClass: string
  onOpenSizePicker: () => void
}

export default function InputNaiParamsPanel({
  cols,
  params,
  setParams,
  displaySize,
  displaySizeLabel,
  selectClass,
  onOpenSizePicker,
}: Props) {
  const [seedDraft, setSeedDraft] = useState(params.seed != null ? String(params.seed) : '')

  const numberInputClass =
    'px-3 py-1.5 rounded-xl border border-gray-200/60 dark:border-white/[0.08] bg-white/50 dark:bg-white/[0.03] focus:outline-none text-xs transition-all duration-200 shadow-sm font-mono'

  return (
    <div className={`grid ${cols} gap-2 text-xs flex-1`}>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">尺寸</span>
        <button
          type="button"
          onClick={() => { dismissAllTooltips(); onOpenSizePicker() }}
          className="px-3 py-1.5 rounded-xl border border-gray-200/60 dark:border-white/[0.08] bg-white/50 dark:bg-white/[0.03] hover:bg-white dark:hover:bg-white/[0.06] focus:outline-none text-xs text-left transition-all duration-200 shadow-sm"
        >
          {displaySizeLabel ?? displaySize}
        </button>
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">步数</span>
        <input
          type="number"
          min={NAI_STEPS_MIN}
          max={NAI_STEPS_MAX}
          value={params.steps ?? NAI_STEPS_DEFAULT}
          onChange={(e) => setParams({ steps: clampNaiSteps(Number(e.target.value)) })}
          className={numberInputClass}
        />
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">Scale</span>
        <input
          type="number"
          min={1}
          max={20}
          step={0.5}
          value={params.scale ?? NAI_SCALE_DEFAULT}
          onChange={(e) => setParams({ scale: clampNaiScale(Number(e.target.value)) })}
          className={numberInputClass}
        />
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">采样器</span>
        <Select
          value={params.sampler ?? 'k_euler_ancestral'}
          onChange={(val) => setParams({ sampler: val })}
          options={NAI_SAMPLER_OPTIONS.map((o) => ({ label: o.label, value: o.value }))}
          className={selectClass}
        />
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">种子</span>
        <input
          type="text"
          inputMode="numeric"
          value={seedDraft}
          onChange={(e) => setSeedDraft(e.target.value.replace(/[^\d]/g, ''))}
          onBlur={() => {
            const parsed = parseNaiSeedInput(seedDraft)
            setParams({ seed: parsed })
            if (parsed == null) setSeedDraft('')
            else setSeedDraft(String(parsed))
          }}
          placeholder="随机"
          className={numberInputClass}
        />
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">噪声表</span>
        <Select
          value={params.noise_schedule ?? 'karras'}
          onChange={(val) => setParams({ noise_schedule: val })}
          options={NAI_NOISE_SCHEDULE_OPTIONS.map((o) => ({ label: o.label, value: o.value }))}
          className={selectClass}
        />
      </label>
      <label className="flex flex-col gap-0.5">
        <span className="text-gray-400 dark:text-gray-500 ml-1">格式</span>
        <Select
          value={params.image_format ?? 'png'}
          onChange={(val) => setParams({ image_format: val as TaskParams['image_format'] })}
          options={[
            { label: 'PNG', value: 'png' },
            { label: 'WebP', value: 'webp' },
          ]}
          className={selectClass}
        />
      </label>
      <label className="flex flex-col gap-0.5 justify-end">
        <span className="text-gray-400 dark:text-gray-500 ml-1">Variety+</span>
        <button
          type="button"
          onClick={() => setParams({ variety_boost: !(params.variety_boost ?? false) })}
          className={`px-3 py-1.5 rounded-xl border text-xs transition-all duration-200 shadow-sm ${
            params.variety_boost
              ? 'border-blue-400 bg-blue-50 text-blue-600 dark:border-blue-500/50 dark:bg-blue-500/10 dark:text-blue-300'
              : 'border-gray-200/60 dark:border-white/[0.08] bg-white/50 dark:bg-white/[0.03] text-gray-600 dark:text-gray-300'
          }`}
        >
          {params.variety_boost ? '开' : '关'}
        </button>
      </label>
    </div>
  )
}
