import { createRootRouteWithContext, Link, Outlet, useRouterState } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import Header from '../components/Header'
import { QueryClient } from '@tanstack/react-query'
import type { Auth } from '../utils/auth'
import { Spinner } from '../components/Spinner'


function RouterSpinner() {
  const isLoading = useRouterState({ select: (s) => s.status === 'pending' })
  return <Spinner show={isLoading} />
}

export const Route = createRootRouteWithContext<{
  auth: Auth
  queryClient: QueryClient
}>()({
  component: RootComponent,
  notFoundComponent: () => {
    return (
      <div>
        <p>This is the notFoundComponent configured on root route</p>
        <Link to="/">Start Over</Link>
      </div>
    )
  },
})
function RootComponent() {
  return (
    <>
      <div>
        < Header />
      </div>
      < RouterSpinner />
      <hr />
      <Outlet />
      <TanStackRouterDevtools position="bottom-left" />
      <ReactQueryDevtools position="right" />
    </>
  )
}