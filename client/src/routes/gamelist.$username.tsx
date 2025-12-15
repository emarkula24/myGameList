import { createFileRoute, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { gameListQueryOptions } from '../queryOptions'
import { useQuery, useQueryErrorResetBoundary, useSuspenseQuery } from '@tanstack/react-query'
import { useState } from 'react'
import { useAuth } from '../utils/auth'
import GameListFilterHeader from '../components/GameListFilterHeader'
import CommonDivider from '../components/CommonDivider'
import { useFilteredGameList } from '../hooks/useGameListFilter'
import { GameListStatusHeader } from '../components/GameListStatusHeader'
import GameListContainer from '../components/GameListContainer'
import { GameListTable } from '../components/GameListTable'
import React from 'react'
import { GameListEmptyError, GameListNotFoundError } from '../game'
export const Route = createFileRoute('/gamelist/$username')({
  loader: ({ context: { queryClient }, params: { username } }) => {
    return queryClient.ensureQueryData(gameListQueryOptions(username))
  },
  component: GameListComponent,
  errorComponent: GameListErrorComponent,
})

function GameListErrorComponent({error}: ErrorComponentProps) {
  const router = useRouter()


  const queryErrorResetBoundary = useQueryErrorResetBoundary()

  React.useEffect(() => {
    queryErrorResetBoundary.reset()
  }, [queryErrorResetBoundary])
  if (error instanceof GameListEmptyError) {
    return <div className='routeContainer'><div style={{ textAlign: "center", fontSize: "2.5em", marginTop: "3%", height: "100%", width: "75%"}}>{error.message}</div></div>
  }
  if (error instanceof GameListNotFoundError) {
    return <div className='routeContainer'><div style={{ textAlign: "center", fontSize: "2.5em", marginTop: "3%", height: "100%", width: "75%"}}>{error.message}</div></div>
  }
  return (
    <div className='routeContainer'>
      <button
        onClick={() => {
          void router.invalidate()
        }}
        type="button"
      >
        retry
      </button>
    </div>
  )
}

function GameListComponent() {

  const statusOptions: Record<number, string> = {
    0: "ALL GAMES", 1: "PLAYING", 2: "COMPLETED", 3: "ON-HOLD", 4: "DROPPED", 5: "PLAN TO PLAY"
  }

  const username = Route.useParams().username
  const auth = useAuth()
  const loggedInUsername = auth.user?.username
  const { data: gamelist } = useSuspenseQuery(gameListQueryOptions(username))
  const { data: loggedInUserGameList = [] } = useQuery({
    ...gameListQueryOptions(loggedInUsername),
    enabled: Boolean(loggedInUsername)
  })
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