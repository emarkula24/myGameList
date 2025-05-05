import { createFileRoute } from '@tanstack/react-router'
export const Route = createFileRoute('/')({
  component: Index,
})

function Index() {
  return (
    <div style={{ padding: "2px" }}>
      <h3>Welcome to myGameList!</h3>

    </div>
  )
}