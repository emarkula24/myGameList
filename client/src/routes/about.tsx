import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/about')({
        component: About,
})

function About() {
  return (
    <div className="p-4 bg-blue-100 text-blue-900 rounded-lg shadow">
      Hello from About!
    </div>
  );
}