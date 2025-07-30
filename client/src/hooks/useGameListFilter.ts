import { useMemo } from "react"
import type { GameListEntries } from "../types/types"

export function useFilteredGameList(
  rawGames: GameListEntries[],
  selectedFilter: number,
  searchQuery: string
): GameListEntries[] {
  return useMemo(() => {
    let filtered = rawGames
    if (selectedFilter !== 0) {
      filtered = filtered.filter((g) => Number(g.status) === selectedFilter)
    }
    filtered = filtered.sort((a, b) => a.status - b.status || a.name.localeCompare(b.name))
    if (searchQuery) {
      filtered = filtered.filter((g) => g.name.toLowerCase().includes(searchQuery.toLowerCase()))
    }
    return filtered
  }, [rawGames, selectedFilter, searchQuery])
}
