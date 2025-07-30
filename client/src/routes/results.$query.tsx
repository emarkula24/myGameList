import { createFileRoute } from '@tanstack/react-router'
import GameResultRow from '../components/GameResultRow'
import { gamesQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/results/$query')({
  loader: ({context: {queryClient}, params: {query}}) => {
    return queryClient.ensureQueryData(gamesQueryOptions(query))
  },
  component: ResultsComponent,
})

function ResultsComponent() {
  const query = Route.useParams().query
  const {data: queriedGames} = useSuspenseQuery(gamesQueryOptions(query))
  return (
    <div className="routeContainer">
      {queriedGames.length > 0 ? (
        queriedGames.map((game) => (
          <GameResultRow game={game} key={game.id}/>
        ))
      ) : (
        <p>No results found.</p>
      )}
    </div>
  )
}