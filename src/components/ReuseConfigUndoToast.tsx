import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useRef } from 'react'
import { useStore } from '../store'

export default function ReuseConfigUndoToast() {
  const reusedConfig = useStore((s) => s.reusedConfig)
  const setReusedConfig = useStore((s) => s.setReusedConfig)
  const undoReuseConfig = useStore((s) => s.undoReuseConfig)
  const undoTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  useEffect(() => {
    if (reusedConfig) {
      // 5秒后自动隐藏
      if (undoTimerRef.current) {
        clearTimeout(undoTimerRef.current)
      }
      undoTimerRef.current = setTimeout(() => {
        setReusedConfig(null)
      }, 5000)
    }

    return () => {
      if (undoTimerRef.current) {
        clearTimeout(undoTimerRef.current)
      }
    }
  }, [reusedConfig, setReusedConfig])

  const handleUndo = () => {
    if (undoTimerRef.current) {
      clearTimeout(undoTimerRef.current)
    }
    undoReuseConfig()
  }

  const handleDismiss = () => {
    if (undoTimerRef.current) {
      clearTimeout(undoTimerRef.current)
    }
    setReusedConfig(null)
  }

  return (
    <AnimatePresence>
      {reusedConfig && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.25, ease: 'easeOut' }}
          className="fixed top-24 left-1/2 -translate-x-1/2 z-50 flex items-center justify-between gap-3 rounded-xl border border-blue-300/70 bg-blue-50/95 px-4 py-2.5 text-sm shadow-lg backdrop-blur-sm dark:border-blue-700/50 dark:bg-blue-900/90"
          style={{ minWidth: '320px', maxWidth: '500px' }}
        >
          <span className="text-blue-800 dark:text-blue-200 font-medium">已复用配置到输入框</span>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleUndo}
              className="rounded-lg bg-blue-600 px-3 py-1 text-xs font-medium text-white hover:bg-blue-700 transition"
            >
              撤回
            </button>
            <button
              type="button"
              onClick={handleDismiss}
              className="text-lg leading-none text-blue-600 hover:text-blue-800 dark:text-blue-300 dark:hover:text-blue-100 transition"
              title="关闭"
            >
              ×
            </button>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  )
}
