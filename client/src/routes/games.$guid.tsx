import { createFileRoute, ErrorComponent, useRouter, type ErrorComponentProps } from '@tanstack/react-router'
import { gameQueryOptions } from '../queryOptions'
import { useQueryErrorResetBoundary, useSuspenseQuery } from '@tanstack/react-query'
import { GameNotFoundError } from '../game'
import React from 'react'


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
      >
        retry
      </button>
      <ErrorComponent error={error} />
    </div>
  )
}

function GameComponent() {
  const guid = Route.useParams().guid
  const { data: game } = useSuspenseQuery(gameQueryOptions(guid))

  return (
    <>
      <div>{game.name}</div>
      <div>{game.deck}</div>

    </>
  )
}
