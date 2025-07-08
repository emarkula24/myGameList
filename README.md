# myGameList
> [!WARNING]
> This project is in a work in progress state
## Table of contents
- [Overview](#Overview)
- [Technologies](#Technologies)
- [Requirements](#Requirements)
- [Install & Run](Install--Run)
- [Design Philosophy](Design&nbsp;philosophy)
- [Test Suite](Test&nbsp;suite)
## Overview
myGameList is an open source platform built to help gaming enthusiasts to keep track of their gaming backlog and connect with a community while doing so. 
The platform is currently self-hosted on a private server, and a containerized version is provided for local hosting.

***All game data used in the project is provided by [Giantbomb](https://www.giantbomb.com/)***
## Technologies
**Backend:** Go, testcontainers

**Frontend:** TypeScript, React, Vite, Vitest, TanStack Router & Query (formerly known as React Query)

**Server:** Docker, Linux (Debian), nginx

**CI/CD:** GitLab, GitHub Actions

A full list of dependencies is available in:
- `api/go.mod` for the backend
- `client/package.json` for the frontend
## Requirements
In the case of someone wanting to hosting their own instance of the app locally, [Docker Desktop](https://docs.docker.com/desktop/setup/install/windows-install/) is recommended.

If not using Docker, the following are required:

```
Node 24.0.2 or newer
Go 1.24.2 or newer
npm serve or similar
```

## Install & Run
**Using Docker**
1. Clone the repository
2. Set up environment variables according to .env.example files
3. Navigate to `/local/` and run:
 `docker compose up -d`
 
4. Application runs `localhost:3004` by default

## Test Suite
A comprehensive test suite has not been implemented yet.

Integration tests for the backend can be found in `api/integration_test`, unit tests share folders with their UUT  (Unit Under Test) and end in `_test.go`.

Frontend testing will begin shortly, testing will be conducted using Vitest and Playwright.
## Design Philosophy
The main goal of this project was to improve as a software developer.

I wanted to create something useful and do deployment, pipelines, and architecture from scratch while properly testing all the application features. The main concern when designing the application
was building it in scalable and testable way with little code repeated. 

On the backend, TDD development style was adopted after figuring out the architecture and Go way of coding. 
For CI/CD pipelines, I decided to use GitLab because the workflow was compatible with my self hosted deployment solution.

