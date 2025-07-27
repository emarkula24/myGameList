import { createFileRoute } from '@tanstack/react-router'
import { useSearch } from '../hooks/useSearchContext'
import GameResultRow from '../components/GameResultRow'

export const Route = createFileRoute('/results')({
  component: RouteComponent,
})

function RouteComponent() {
  const {searchResults} = useSearch()
  console.log(searchResults)
  return (
    <div className="routeContainer">
      {searchResults.length > 0 ? (
        searchResults.map((game) => (
          <GameResultRow game={game} />
        ))
      ) : (
        <p>No results found.</p>
      )}
    </div>
  )
}