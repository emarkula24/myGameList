import { createContext, use } from 'react'
import type { Games } from '../types/types'

type SearchContextType = {
  searchResults: Games[]
  setSearchResults: React.Dispatch<React.SetStateAction<Games[]>>
}

export const SearchContext = createContext<SearchContextType | undefined>(undefined)

export const useSearch = () => {
  const ctx = use(SearchContext)
  if (!ctx) throw new Error('useSearch must be used within SearchProvider')
  return ctx
}