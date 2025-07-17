
import { createFileRoute, Link } from '@tanstack/react-router'
import { auth } from '../utils/auth'

export const Route = createFileRoute('/')({
        component: IndexComponent,
})

function IndexComponent() {
        
    return (
        <div style={{ padding: "2px" }}>
            <h3>Welcome to myGameList!</h3>

            {auth.status === 'loggedIn' && auth.username && (
                <Link to="/gamelist/$username" params={{ username: auth.username }}>
                    View your gamelist!
                </Link>
            )}
        </div>
    )
}