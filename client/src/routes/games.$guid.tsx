import { createFileRoute, ErrorComponent, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { gameQueryOptions } from '../queryOptions'
import { useMutation, useQueryErrorResetBoundary, useSuspenseQuery } from '@tanstack/react-query'
import { addGame, GameNotFoundError } from '../game'
import React from 'react'
import { useAuth } from '../utils/auth'
import GameAddDropdown from '../components/GameAddDropdown'


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
  const router = useRouter()
  const auth = useAuth()
  const guid = Route.useParams().guid
  const { data: game } = useSuspenseQuery(gameQueryOptions(guid))

  const addMutation = useMutation({
    mutationFn: async (status: string) => {
      return await addGame(game.id, status, auth.user?.username, game.name)
    },
    onSuccess: () => {
      router.invalidate()
    },
    onError: (error) => {
      console.error(error)
    }
  })

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

      {auth.isAuthenticated && (
        <GameAddDropdown
          onSelect={(status) => addMutation.mutate(status)}
        />
      )}
    </div>
  )
}
