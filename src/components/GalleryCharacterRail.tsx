import { useStore } from '../store'
import CharacterRail, { getCharacterGridItems } from './input/CharacterRail'
import PositionGrid from './input/PositionGrid'

export default function GalleryCharacterRail() {
  const characters = useStore((s) => s.characters)
  const useCoords = useStore((s) => s.use_coords)
  const useOrder = useStore((s) => s.use_order)
  const addCharacter = useStore((s) => s.addCharacter)
  const updateCharacter = useStore((s) => s.updateCharacter)
  const removeCharacter = useStore((s) => s.removeCharacter)
  const setUseCoords = useStore((s) => s.setUseCoords)
  const setUseOrder = useStore((s) => s.setUseOrder)

  return (
    <aside className="hidden lg:block w-[300px] shrink-0 sticky top-24 self-start max-h-[calc(100vh-7rem)] overflow-hidden rounded-3xl border border-gray-300/70 bg-gray-100/82 p-4 shadow-[0_18px_50px_rgba(15,23,42,0.14)] ring-1 ring-black/10 backdrop-blur-2xl supports-[backdrop-filter]:bg-gray-100/72 dark:border-white/[0.10] dark:bg-gray-900/78 dark:shadow-[0_18px_50px_rgba(0,0,0,0.35)] dark:ring-white/10">
      <CharacterRail
        characters={characters}
        useCoords={useCoords}
        useOrder={useOrder}
        onAddCharacter={addCharacter}
        onUpdateCharacter={updateCharacter}
        onRemoveCharacter={removeCharacter}
        onUseCoordsChange={setUseCoords}
        onUseOrderChange={setUseOrder}
      />
    </aside>
  )
}

export function GalleryCharacterOverview() {
  const characters = useStore((s) => s.characters)
  const useCoords = useStore((s) => s.use_coords)
  const gridCharacters = getCharacterGridItems(characters)

  return (
    <aside className="hidden xl:block w-[230px] shrink-0 sticky top-24 self-start rounded-3xl border border-gray-300/70 bg-gray-100/82 p-4 shadow-[0_18px_50px_rgba(15,23,42,0.14)] ring-1 ring-black/10 backdrop-blur-2xl supports-[backdrop-filter]:bg-gray-100/72 dark:border-white/[0.10] dark:bg-gray-900/78 dark:shadow-[0_18px_50px_rgba(0,0,0,0.35)] dark:ring-white/10">
      <div className="mb-3 flex items-center justify-between gap-2">
        <span className="text-xs font-semibold tracking-wide text-gray-700 dark:text-gray-200">总览</span>
        <span className="rounded-full border border-gray-300/60 bg-white/45 px-2 py-0.5 text-[10px] text-gray-500 dark:border-white/[0.08] dark:bg-white/[0.04] dark:text-gray-400">{characters.length} 人</span>
      </div>
      {useCoords && characters.length > 0 ? (
        <PositionGrid
          value={characters[0]?.position ?? 'C3'}
          onChange={() => {}}
          characters={gridCharacters}
          className="mx-auto max-w-[190px]"
        />
      ) : (
        <div className="rounded-2xl border border-dashed border-gray-300/80 bg-white/35 px-3 py-8 text-center text-[11px] text-gray-500 shadow-inner dark:border-white/[0.08] dark:bg-white/[0.03] dark:text-gray-500">
          {characters.length === 0 ? '添加角色后显示方位总览' : '开启 Manual coords 后显示方位总览'}
        </div>
      )}
    </aside>
  )
}
