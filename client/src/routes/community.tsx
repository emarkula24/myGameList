import { createFileRoute, Link } from '@tanstack/react-router'
import CommonDivider from '../components/CommonDivider'
import { usersQueryOptions } from '../queryOptions'
import { useSuspenseQuery } from '@tanstack/react-query'

export const Route = createFileRoute('/community')({
  loader: ({ context: {queryClient}}) => {
    return queryClient.ensureQueryData(usersQueryOptions())
  },
  component: RouteComponent,
})

function RouteComponent() {
  const {data: users} = useSuspenseQuery(usersQueryOptions()) 
  console.log(users)
  return (
    <div className="routeContainer">
      <CommonDivider routeName={"Community"} />
      <h1>See other people's lists through here!</h1>
      <div>
        {users.map((user) => (
          <Link 
          key={user.id}
          to="/gamelist/$username"
          params={{username: user.username}}
          >
            {user.username} <br></br>
            </Link>
        ))}
      </div>
    </div>
  )

}
