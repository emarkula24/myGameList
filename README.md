# myGameList
> [!WARNING]
> This project is in a work in progress state
## Table of contents
- [Overview](Overview)
- [Technologies](Technologies)
- [Requirements](Requirements)
- [Install & Run](Install--Run)
- [Design Philosophy](Design&nbsp;philosophy)
- [Test Suite](Test&nbsp;suite)
## Overview
myGameList is an open source platform built to help gaming enthusiasts to keep track of their gaming backlog and connect with a community while doing so. 
The platform is currently self-hosted on a private linux server, and a containerized version is provided for local hosting.

***All game data used in the project is provided by [Giantbomb](https://www.giantbomb.com/)***
## Technologies
**Backend:** Go, testcontainers

**Frontend:** TypeScript, React, Vite, Vitest, PlayWright, TanStack Router & Query (formerly known as React Query)

**Server:** Docker, Linux (Debian), nginx

**CI/CD:** GitLab, GitHub Actions

A full list of dependencies is available in:
- `api/go.mod` for the backend
- `client/package.json` for the frontend
## Requirements
In the case of someone wanting to hosting their own instance of the app locally, [Docker Desktop](https://docs.docker.com/desktop/setup/install/windows-install/) is recommended.

## Install & Run
**Using Docker**
1. Clone the repository
2. Set up environment variables according to .env.example files
3. Navigate to `/local/` and run:
 `docker compose up -d`
 
4. Application runs `localhost:3004` by default

## Test Suite
Integration tests for the backend can be found in `api/integration_test`, unit tests share folders with their UUT  (Unit Under Test) and end in `_test.go`.

Frontend testing will begin shortly, testing will be conducted using Vitest and Playwright.
## Design Philosophy
The main goal of this project was to improve as a software developer.
I wanted to create something useful and do deployment, pipelines, and architecture from scratch while properly testing all the application features. There was no one to say to me what is right and what is wrong, so I tried to make a working product as scalable, maintainable, and with as little repeated code as possible in my own ways.

On the backend, TDD development style was adopted after figuring out the architecture and Go way of coding. I decided to not implement many unit tests because of the simplicity of the API endpoints. External APIs are mocked in tests to account for negative scenarios. For CI/CD pipelines, I decided to use GitLab because the workflow was compatible with my self hosted deployment solution. 

