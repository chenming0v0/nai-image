import { AnimatePresence, motion } from 'framer-motion'
import { useEffect, useRef } from 'react'
import { useStore } from '../store'

export default function CharacterDeleteUndoToast() {
  const deletedCharacter = useStore((s) => s.deletedCharacter)
  const setDeletedCharacter = useStore((s) => s.setDeletedCharacter)
  const undoDeleteCharacter = useStore((s) => s.undoDeleteCharacter)
  const undoTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)

  useEffect(() => {
    if (deletedCharacter) {
      // 5秒后自动隐藏
      if (undoTimerRef.current) {
        clearTimeout(undoTimerRef.current)
      }
      undoTimerRef.current = setTimeout(() => {
        setDeletedCharacter(null)
      }, 5000)
    }

    return () => {
      if (undoTimerRef.current) {
        clearTimeout(undoTimerRef.current)
      }
    }
  }, [deletedCharacter, setDeletedCharacter])

  const handleUndo = () => {
    if (undoTimerRef.current) {
      clearTimeout(undoTimerRef.current)
    }
    undoDeleteCharacter()
  }

  const handleDismiss = () => {
    if (undoTimerRef.current) {
      clearTimeout(undoTimerRef.current)
    }
    setDeletedCharacter(null)
  }

  return (
    <AnimatePresence>
      {deletedCharacter && (
        <motion.div
          initial={{ opacity: 0, y: -20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.25, ease: 'easeOut' }}
          className="fixed top-24 left-1/2 -translate-x-1/2 z-50 flex items-center justify-between gap-3 rounded-xl border border-amber-300/70 bg-amber-50/95 px-4 py-2.5 text-sm shadow-lg backdrop-blur-sm dark:border-amber-700/50 dark:bg-amber-900/90"
          style={{ minWidth: '320px', maxWidth: '500px' }}
        >
          <span className="text-amber-800 dark:text-amber-200 font-medium">已删除角色</span>
          <div className="flex items-center gap-2">
            <button
              type="button"
              onClick={handleUndo}
              className="rounded-lg bg-amber-600 px-3 py-1 text-xs font-medium text-white hover:bg-amber-700 transition"
            >
              撤回
            </button>
            <button
              type="button"
              onClick={handleDismiss}
              className="text-lg leading-none text-amber-600 hover:text-amber-800 dark:text-amber-300 dark:hover:text-amber-100 transition"
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
