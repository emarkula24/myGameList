import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '../utils/auth'

export const Route = createFileRoute('/_auth/profile')({
  component: ProfileComponent,
})

function ProfileComponent() {
  const auth = useAuth()
  return (
    <div className="routeContainer">
      <h1>Username:<strong>{auth.user?.username}</strong></h1>
    </div>
  )
}