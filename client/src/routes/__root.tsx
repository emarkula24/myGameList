import { createRootRouteWithContext, Link, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'
import Header from '../components/Header'
import { QueryClient } from '@tanstack/react-query'
import type { AuthContext } from '../utils/auth'
import { useMemo, useState } from 'react'
import type { Games } from '../types/types'
import { SearchContext } from '../hooks/useSearchContext'
import MainHeader from '../components/MainHeader'


// function RouterSpinner() {
//   const isLoading = useRouterState({ select: (s) => s.status === 'pending' })
//   return <Spinner show={isLoading} />
// }
interface MyRouterContext {
  auth: AuthContext
  queryClient: QueryClient
}
export const Route = createRootRouteWithContext<MyRouterContext>()({

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
  const [searchResults, setSearchResults] = useState<Games[]>([])
  const searchContextValue = useMemo(() => ({ searchResults, setSearchResults }), [searchResults])
  return (
    <div >
      <SearchContext value={searchContextValue}>
        
        <div style={{display:"flex", flexDirection: "column", alignItems: "center", justifyContent: "center"}}>
          < MainHeader />
          < Header />

        </div>
        <hr />
        <Outlet />
        <TanStackRouterDevtools position="bottom-left" />
        <ReactQueryDevtools position="right" />
      </SearchContext>
    </div>
  )
}