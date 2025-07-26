import { createFileRoute } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useState } from 'react'
import GameRow from '../components/GameRow'
import { useAuth } from '../utils/auth'
import GameListFilterHeader from '../components/GameListFilterHeader'
import React from 'react'
import styles from "./gamelist.$username.module.css"
export const Route = createFileRoute('/gamelist/$username')({
  loader: ({ context: { queryClient }, params: { username } }) => {
    return queryClient.ensureQueryData(gameListQueryOptions(username))
  },
  component: GameListComponent,
})

function GameListComponent() {

  const statusOptions: { [key: number]: string } = {
    0: "ALL GAMES",
    1: "PLAYING",
    2: "COMPLETED",
    3: "ON-HOLD",
    4: "DROPPED",
    5: "PLAN TO PLAY"
  }

  const username = Route.useParams().username
  const auth = useAuth()
  const { data: gamelist } = useSuspenseQuery(gameListQueryOptions(username))
  const { data: loggedInUserGameList } = useSuspenseQuery(gameListQueryOptions(auth.user?.username))
  const [selectedFilter, setSelectedFilter] = useState(0)
  const [editingGameIds, setEditingGameIds] = useState<Set<number>>(new Set())
  const loggedInGameIds = new Set(loggedInUserGameList.map(g => g.id))

  const startEditing = (gameId: number) => {
    setEditingGameIds(prev => {
      const newSet = new Set(prev)
      newSet.add(gameId)
      return newSet
    });
  };
  const stopEditing = (gameId: number) => {
    setEditingGameIds(prev => {
      const newSet = new Set(prev)
      newSet.delete(gameId)
      return newSet
    })
  }
  React.useEffect(() => {
    setEditingGameIds(new Set())
  }, [selectedFilter])

  let filteredGameList = gamelist
  if (selectedFilter !== 0) {
    filteredGameList = gamelist.filter((game) => Number(game.status) === selectedFilter)
  }
  filteredGameList.sort((a, b) => {
    //compare status numerically
    if (a.status !== b.status) {
      return a.status - b.status
    }
    // If status is equal, compare names alphabetically
    return a.name.localeCompare(b.name)
  })
  return (
    <div>
      <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
        <GameListFilterHeader onSelect={setSelectedFilter} />
        <div style={{ padding: "8px" }}></div>
        <div className={styles.divider}>
          {/* <span>Viewing {username}'s Game List</span> */}
          <span className={styles.text}>{statusOptions[selectedFilter]}</span>
        </div>
        <div style={{ display: "flex", flexDirection: "column", width: "1000px", border: "solid lightgrey 1px"}}>
          {filteredGameList.length > 0 ? (
            <table>
              <thead>
                <tr style={{backgroundColor: "lightgray"}}>
                  {/* <th className={styles.headerTitleStatus}></th> */}
                  <th className={styles.headerNumber}>#</th>
                  <th className={styles.headerImage}>Image</th>
                  <th className={styles.headerName}>Game Title</th>
                  <th className={styles.headerStatus}>Status</th>
                  <th className={styles.headerAction}>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredGameList.map((game, index) => (
                  <GameRow
                    index={index}
                    key={game.id}
                    game={game}
                    username={username}
                    isEditing={editingGameIds.has(game.id)}
                    startEditing={startEditing}
                    stopEditing={stopEditing}
                    onUpdateSuccess={() => stopEditing(game.id)}
                    isMissingFromLoggedInUserList={!loggedInGameIds.has(game.id)}
                  />
                ))}
              </tbody>
            </table>
            
          ) : (
            <p>This category is empty</p>
          )}
        </div>
      </div>
    </div>
  )
}
