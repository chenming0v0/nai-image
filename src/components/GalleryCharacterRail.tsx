import { useStore } from '../store'
import CharacterRail from './input/CharacterRail'

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
