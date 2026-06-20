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
    <aside className="hidden lg:block w-[300px] shrink-0 sticky top-24 self-start max-h-[calc(100vh-7rem)] rounded-3xl border border-white/50 bg-white/70 p-4 shadow-[0_8px_30px_rgb(0,0,0,0.06)] ring-1 ring-black/5 backdrop-blur-2xl dark:border-white/[0.08] dark:bg-gray-900/70 dark:shadow-[0_8px_30px_rgb(0,0,0,0.25)] dark:ring-white/10">
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
