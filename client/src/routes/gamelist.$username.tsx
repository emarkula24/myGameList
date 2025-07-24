import { createFileRoute, Link } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/gamelist/$username')({
    loader: ({ context: {queryClient}, params: {username} }) => {
        return queryClient.ensureQueryData(gameListQueryOptions(username))
},
  component: GameListComponent,
})

const statusOptions: { [key: number]: string } = {
    1: "Playing",
    2: "Completed",
    3: "On-Hold",
    4: "Dropped",
    5: "Plan to Play"
}

function GameListComponent() {
  const username = Route.useParams().username
  const {data: gamelist} = useSuspenseQuery(gameListQueryOptions(username))
  console.log(gamelist)
  gamelist.sort((a, b) => a.name.localeCompare(b.name));
  return (
    <div>
      <h3>Viewing {username}'s Game List</h3>
      {gamelist.length > 0 ? (
        gamelist.map((game) => (
          <div key={game.id}>
            <table>
              <thead>
                <tr>
                  <th>Image</th>
                  <th>Game Title</th>
                  <th>Status</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <th><img src={game.image.icon_url} alt="" /></th>
                  <th><Link to={"/games/$guid"} params={{guid: game.guid}}>{game.name}</Link></th>
                  <th>{statusOptions[game.status]}</th>
                </tr>
              </tbody>
            </table>
          </div>
          
        ))
      ) : (
        <p>No results found.</p>
      )}

    </div>
  )
}

