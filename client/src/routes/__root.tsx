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
import Footer from '../components/Footer'


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
      <div className='routeContainer' style={{ fontSize: "2em" }}>
        <p style={{fontSize: "2em"}}>404 address not found</p>
        <Link to="/">Top Page</Link>
      </div>
    )
  },
})
function RootComponent() {
  const [searchResults, setSearchResults] = useState<Games[]>([])
  const searchContextValue = useMemo(() => ({ searchResults, setSearchResults }), [searchResults])
  return (
    <div>
      <SearchContext value={searchContextValue}>

        <div style={{ display: "flex", flexDirection: "column", alignItems: "center" }}>
          < MainHeader />
          < Header />
        </div>
        <Outlet />
        <div style={{ borderTop: "1px solid lightgrey", width: "100%" }}></div>
        <div style={{ display: 'flex', alignItems: "center", justifyContent: "center" }}>
          <Footer />
        </div>




        <TanStackRouterDevtools position="bottom-left" />
        <ReactQueryDevtools position="right" />
      </SearchContext>

    </div>
  )
}