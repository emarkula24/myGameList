import { createFileRoute } from '@tanstack/react-router'
import { queryClient } from '../utils/queryClient'
import { gameQueryOptions } from './- components/gameQueryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/game/$gameGuid')({
  loader: ({params: {gameGuid}}) => {
    return queryClient.ensureQueryData(gameQueryOptions(gameGuid))
  },
  component: GameComponent,
})

function GameComponent() {
  const gameGuid = Route.useParams().gameGuid
  const { data: game} = useSuspenseQuery(gameQueryOptions(gameGuid))


  
  return <div>Hello "/game/$gameId"!</div>
}
