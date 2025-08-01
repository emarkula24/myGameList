import { createFileRoute } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { useAuth } from '../utils/auth'
import GameListFilterHeader from '../components/GameListFilterHeader'
import CommonDivider from '../components/CommonDivider'
import { useFilteredGameList } from '../hooks/useGameListFilter'
import { GameListStatusHeader } from '../components/GameListStatusHeader'
import GameListContainer from '../components/GameListContainer'
import { GameListTable } from '../components/GameListTable'
export const Route = createFileRoute('/gamelist/$username')({
  loader: ({ context: { queryClient }, params: { username } }) => {
    return queryClient.ensureQueryData(gameListQueryOptions(username))
  },
  component: GameListComponent,
})

function GameListComponent() {

  const statusOptions: Record<number, string> = {
    0: "ALL GAMES", 1: "PLAYING", 2: "COMPLETED", 3: "ON-HOLD", 4: "DROPPED", 5: "PLAN TO PLAY"
  }

  const username = Route.useParams().username
  const auth = useAuth()
  const { data: gamelist } = useSuspenseQuery(gameListQueryOptions(username))
  const { data: loggedInUserGameList } = useSuspenseQuery(gameListQueryOptions(auth.user?.username))
  const [selectedFilter, setSelectedFilter] = useState(0)
  const [searchQuery, setSearchQuery] = useState("")
  const [editingGameIds, setEditingGameIds] = useState<Set<number>>(() => new Set())
  const loggedInGameIds = new Set(loggedInUserGameList.map(g => g.id))
  const filteredGameList = useFilteredGameList(gamelist, selectedFilter, searchQuery)

  const handleFilterChange = (filter: number) => {
    setSelectedFilter(filter)
    setEditingGameIds(new Set())

  }
  return (
    <div className="routeContainer">
      < CommonDivider routeName={`Viewing ${username}'s Game List`} />
      <GameListContainer>
        <GameListFilterHeader
          onSelect={handleFilterChange}
          setSearchQuery={setSearchQuery}
        />
        <div style={{ padding: "8px" }}></div>
        <GameListStatusHeader statusText={statusOptions[selectedFilter]} />
        <div style={{ display: "flex", flexDirection: "column", width: "100%", border: "solid lightgrey 1px", }}>
          <GameListTable
            games={filteredGameList}
            username={username}
            editingGameIds={editingGameIds}
            startEditing={(id) => setEditingGameIds(new Set(editingGameIds).add(id))}
            stopEditing={(id) => {
              const newSet = new Set(editingGameIds)
              newSet.delete(id)
              setEditingGameIds(newSet)
            }}
            loggedInGameIds={loggedInGameIds}
          />
        </div>
      </GameListContainer>
    </div>
  )
}