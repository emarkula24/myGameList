
import { createFileRoute, Link } from '@tanstack/react-router'
import { useAuth } from '../utils/auth'
import CommonDivider from '../components/CommonDivider'
import WelcomeText from '../components/WelcomeText'

export const Route = createFileRoute('/')({
    component: IndexComponent,
})

function IndexComponent() {
    const { user, isAuthenticated } = useAuth()
    return (
        <div className="routeContainer">
            <CommonDivider routeName='Welcome to MyGameList!' />
            {isAuthenticated && user && (
                <Link to="/gamelist/$username" params={{ username: user.username }} style={{ fontSize: "2.4em" }}>
                    View your gamelist!
                </Link>
            )}
            <WelcomeText />
        </div>
    )
}