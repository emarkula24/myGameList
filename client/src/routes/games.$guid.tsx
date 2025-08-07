import { createFileRoute, ErrorComponent, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { gameListEntryQueryOptions, gameQueryOptions, useAddGameMutation, useUpdateGameMutation } from '../queryOptions'
import { useQueryErrorResetBoundary, useSuspenseQuery } from '@tanstack/react-query'
import { GameNotFoundError } from '../game'
import React from 'react'
import { useAuth } from '../utils/auth'
import GameAddButton from '../components/GameAddButton'
import GameUpdateDropdown from '../components/GameUpdateDropdown'
import CommonDivider from '../components/CommonDivider'


export const Route = createFileRoute('/games/$guid')({
  loader: ({ context: { queryClient }, params: { guid } }) => {
    return queryClient.ensureQueryData(gameQueryOptions(guid))
  },
  component: GameComponent,
  errorComponent: gameErrorComponent,
})

function gameErrorComponent({ error }: ErrorComponentProps) {
  const router = useRouter()
  if (error instanceof GameNotFoundError) {
    return <div>{error.message}</div>
  }

  const queryErrorResetBoundary = useQueryErrorResetBoundary()

  React.useEffect(() => {
    queryErrorResetBoundary.reset()
  }, [queryErrorResetBoundary])

  return (
    <div>
      <button
        onClick={() => {
          void router.invalidate()
        }}
        type="button"
      >
        retry
      </button>
      <ErrorComponent error={error} />
    </div>
  )
}

function GameComponent() {
  const auth = useAuth()
  const guid = Route.useParams().guid
  const { data: game } = useSuspenseQuery(gameQueryOptions(guid))
  const gameListEntryQuery = useSuspenseQuery(gameListEntryQueryOptions(auth.user?.username, game.id))
  const gameListEntry = gameListEntryQuery.data
  const status = gameListEntry.gamedata?.status ?? false

  const addMutation = useAddGameMutation(game.id, game.name)
  const updateMutation = useUpdateGameMutation(auth.user?.username, game.id, game.name)

  return (
    <div className="routeContainer">
      <CommonDivider routeName={game.name} />
      <div
        style={{
          width: "75%",
          display: "grid",
          gridTemplateColumns: "1fr 2fr 0.5fr",
          gridTemplateRows: "auto auto",
          gap: "1rem",
          alignItems: "start",
          justifyContent: "center",

        }}
      >
        <img src={game.image?.medium_url} style={{ width: "100%", height: "auto", gridColumn: "1", gridRow: "1 / span 2"  }} />
        <div style={{ gridColumn: "2", gridRow: "1", fontSize: "1em", }}>
        <p>{game.deck}</p>
        {/* <GamePlatformHeader game={game} /> */}

        {addMutation.isError && (
          <div style={{ color: 'red' }}>Error: {addMutation.error.message}</div>
        )}

        {addMutation.isSuccess && (
          <div style={{ color: 'green' }}>Game successfully added!</div>
        )}

        {addMutation.isPending && (
          <div style={{ color: 'gray' }}>Adding game...</div>
        )}
        {updateMutation.isPending && (
          <div style={{ color: "blue" }}>Updating game status...</div>
        )}
        {updateMutation.isSuccess && (
          <div style={{ color: 'green' }}>Game successfully updated!</div>
        )}
        </div>
        <div style={{ gridColumn: "2", gridRow: "2", justifyContent: "flex-end"}}>
        {auth.isAuthenticated && status ? (
          <GameUpdateDropdown
            onUpdateListEntry={(status) => updateMutation.mutate(status)}
            status={status}
          />
          
        ) : (
          !addMutation.isSuccess && (
            <GameAddButton
              onNewListEntry={(status) => addMutation.mutate(status)}
            />
          )
        )}
        </div>
      </div>
    </div>
  )
}
