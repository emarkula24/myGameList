import { createFileRoute } from '@tanstack/react-router'
import CommonDivider from '../components/CommonDivider'

export const Route = createFileRoute('/community')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className="routeContainer">
      <CommonDivider routeName={"Community"} />
      <h1>meet new friends here</h1>
    </div>
  )

}
