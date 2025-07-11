import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { ErrorComponent, RouterProvider, createRouter } from '@tanstack/react-router'
// Import the generated route tree
import { routeTree } from './routeTree.gen'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { Spinner } from './components/Spinner'
import { auth } from './utils/auth'



const queryClient = new QueryClient()
// Set up a Router instance
const router = createRouter({
        routeTree,
        defaultPendingComponent: () => (
                <div>
                        <Spinner />
                </div>
        ),
        defaultErrorComponent: ({ error }) => <ErrorComponent error={error} />,
        context: {
                auth: undefined!, // We'll inject this when we render
                queryClient,
        },
        defaultPreload: 'intent',
        // Since we're using React Query, we don't want loader calls to ever be stale
        // This will ensure that the loader is always called when the route is preloaded or visited
        defaultPreloadStaleTime: 0,
        scrollRestoration: true,
})
// Register the router instance for type safety
declare module '@tanstack/react-router' {
        interface Register {
                router: typeof router
        }
}

// Render the app
const rootElement = document.getElementById('root')!
if (!rootElement.innerHTML) {
        const root = ReactDOM.createRoot(rootElement)
        root.render(
                <StrictMode>
                        <QueryClientProvider client={queryClient}>
                                <RouterProvider router={router} context={{
                                        auth, 
                                }}
                                />
                        </QueryClientProvider>
                </StrictMode>,
        )
}