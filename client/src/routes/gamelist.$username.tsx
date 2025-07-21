import { createFileRoute, Link } from '@tanstack/react-router'
import { fetchGameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/gamelist/$username')({
    loader: ({ context: {queryClient}, params: {username} }) => {
        return queryClient.ensureQueryData(fetchGameListQueryOptions(username))
},
  component: GameListComponent,
})
function GameListComponent() {
  const username = Route.useParams().username
  const {data: gamelist} = useSuspenseQuery(fetchGameListQueryOptions(username))
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
                  <th>{game.status}</th>
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

