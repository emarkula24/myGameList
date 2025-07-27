import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/about')({
        component: About,
})

function About() {
  return (
    <div className="routeContainer">
      <h1>Hello from About!</h1>
    </div>
  );
}