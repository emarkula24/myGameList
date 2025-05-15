import { createFileRoute } from '@tanstack/react-router'
import { queryClient } from '../utils/queryClient'
import { gameQueryOptions } from './- components/gameQueryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'
export const Route = createFileRoute('/game/$guid')({
  loader: ({params: {guid}}) => {
    return queryClient.ensureQueryData(gameQueryOptions(guid))
  },
  component: GameComponent,
})

function GameComponent() {
  const guid = Route.useParams().guid
  const { data: game} = useSuspenseQuery(gameQueryOptions(guid))

  

  
  return (
    <>
      <div>{game.name}</div>
      <div>{game.deck}</div>
        
    </>
  )
}
