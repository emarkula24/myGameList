import { createFileRoute } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useState } from 'react'
import GameRow from '../components/GameRow'
import { useAuth } from '../utils/auth'

export const Route = createFileRoute('/gamelist/$username')({
  loader: ({ context: { queryClient }, params: { username } }) => {
    return queryClient.ensureQueryData(gameListQueryOptions(username))
  },
  component: GameListComponent,
})

function GameListComponent() {
  const username = Route.useParams().username
  const auth = useAuth()
  const { data: gamelist } = useSuspenseQuery(gameListQueryOptions(username))
  const { data: loggedInUserGameList } = useSuspenseQuery(gameListQueryOptions(auth.user?.username))
  const [editingGameIds, setEditingGameIds] = useState<Set<number>>(new Set())
  const loggedInGameIds = new Set(loggedInUserGameList.map(g => g.id))
  const toggleEditMode = (gameId: number) => {
    setEditingGameIds(prev => {
      const newSet = new Set(prev)
      if (newSet.has(gameId)) {
        newSet.delete(gameId)
      } else {
        newSet.add(gameId)
      }
      return newSet
    })
  }
  gamelist.sort((a, b) => a.name.localeCompare(b.name));
  return (
    <div>
      <h3>Viewing {username}'s Game List</h3>
      {gamelist.length > 0 ? (  
              <table>
                <thead>
                  <tr>
                    <th>Image</th>
                    <th>Game Title</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody>
                  {gamelist.map((game) => (
                    <GameRow
                      key={game.id}
                      game={game}
                      username={username}
                      isEditing={editingGameIds.has(game.id)}
                      toggleEditMode={toggleEditMode}
                    />
                  ))}
           </tbody>
        </table>
      ) : (
        <p>The GameList is empty</p>
      )}
    </div>
  )
}
