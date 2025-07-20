import { createFileRoute } from '@tanstack/react-router'
import { useAuth } from '../utils/auth'

export const Route = createFileRoute('/_auth/profile')({
  component: ProfileComponent,
})

function ProfileComponent() {
  const auth = useAuth()
  return (
    <div className="p-2 space-y-2">
      <div>
        Username:<strong>{auth.user?.username}</strong>
      </div>
    </div>
  )
}