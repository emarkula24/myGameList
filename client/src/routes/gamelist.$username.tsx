import { createFileRoute } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/gamelist/$username')({
    loader: ({ context: {queryClient}, params: {username} }) => {
        return queryClient.ensureQueryData(gameListQueryOptions(username))
},
  component: GameListComponent,
})

function GameListComponent() {
  const username = Route.useParams().username
  const {data: gamelist} = useSuspenseQuery(gameListQueryOptions(username))
  return (
    <div>
      <h3>Viewing {username}'s Game List</h3>
      {gamelist.length > 0 ? (
        gamelist.map((game, index) => (
          <div key={index}>
            <p>{game.id}</p>
            <h3>{game.status}</h3>
          </div>
        ))
      ) : (
        <p>No results found.</p>
      )}

    </div>
  )
}

