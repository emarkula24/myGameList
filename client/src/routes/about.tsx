import { createFileRoute } from '@tanstack/react-router'
import CommonDivider from '../components/CommonDivider'

export const Route = createFileRoute('/about')({
  component: About,
})

function About() {
  return (
    <div className='routeContainer'>
      <CommonDivider routeName='About' />
      <div style={{ width: "75%" }}>

        <h1>About this site</h1>
        <p style={{marginBottom: "1em"}}>
          I created this site as a small project in order to improve my software development skills.
          UI design was kept simple because designing a good UI is a lot harder than coding one.
        </p>
        <p>
          The backend code was built with Go and the frontend with TypeScript and React.js.
          Normal CSS was used instead of a framework like tailwindcss or Next.js because I wanted to practice using it.<br />
          The website is currently self-hosted on my own Linux server, containerized into Docker containers.
        </p>
        <p>
          Please check the GitHub link on the bottom header if you are interested in the details.
        </p>

      </div>
    </div>
  );
}