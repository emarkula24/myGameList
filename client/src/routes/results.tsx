import { createFileRoute } from '@tanstack/react-router'
import { useSearch } from '../hooks/useSearchContext'

export const Route = createFileRoute('/results')({
  component: RouteComponent,
})

function RouteComponent() {
  const {searchResults} = useSearch()
  return (
    <div>
      {searchResults.length > 0 ? (
        searchResults.map((game, index) => (
          <div key={index}>
            <h3>{game.name}</h3>
            <img src={game.image?.thumb_url} />
          </div>
        ))
      ) : (
        <p>No results found.</p>
      )}
    </div>
  )
}