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
    <aside className="hidden lg:flex w-[300px] flex-shrink-0 h-full pr-6 border-r border-gray-300/50 dark:border-white/[0.08]">
      <div className="w-full h-full overflow-y-auto hide-scrollbar">
        <div className="pb-48">
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
        </div>
      </div>
    </aside>
  )
}

export function GalleryCharacterOverview() {
  const characters = useStore((s) => s.characters)
  const useCoords = useStore((s) => s.use_coords)
  const gridCharacters = getCharacterGridItems(characters)

  return (
    <aside className="hidden xl:flex w-[230px] flex-shrink-0 h-full pl-6 border-l border-gray-300/50 dark:border-white/[0.08]">
      <div className="w-full h-full overflow-y-auto hide-scrollbar">
        <div className="pb-48">
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
        </div>
      </div>
    </aside>
  )
}
