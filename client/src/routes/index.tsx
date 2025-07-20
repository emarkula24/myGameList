
    import { createFileRoute, Link } from '@tanstack/react-router'
    import { useAuth } from '../utils/auth'

    export const Route = createFileRoute('/')({
            component: IndexComponent,
    })
    
    function IndexComponent() {
        const { user, isAuthenticated } = useAuth()    
        return (
            <div style={{ padding: "2px" }}>
                <h3>Welcome to myGameList!</h3>

                {isAuthenticated && user && (
                    <Link to="/gamelist/$username" params={{ username: user.username }}>
                        View your gamelist!
                    </Link>
                )}
            </div>
        )
    }