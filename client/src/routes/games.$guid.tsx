  import { createFileRoute, ErrorComponent, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
  import { gameListEntryQueryOptions, gameQueryOptions, useAddGameMutation, useUpdateGameMutation } from '../queryOptions'
  import { useQueryErrorResetBoundary, useSuspenseQuery } from '@tanstack/react-query'
  import { GameNotFoundError } from '../game'
  import React from 'react'
  import { useAuth } from '../utils/auth'
import GameAddButton from '../components/GameAddButton'
import GameUpdateDropdown from '../components/GameUpdateDropdown'


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
            router.invalidate()
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
      <div>
        <h2>{game.name}</h2>
        <p>{game.deck}</p>

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
          <div style={{color: "blue" }}>Updating game status...</div>
        )}
        {updateMutation.isSuccess && (
          <div style={{ color: 'green' }}>Game successfully updated!</div>
        )}
        
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
    )
  }
