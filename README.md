# myGameList
> [!NOTE]
> This project is in a work in progress state
## Table of contents
- [Overview](Overview)
- [Technologies](Technologies)
- [Requirements](Requirements)
- [Install & Run](Install--Run)
- [Design Philosophy](Design&nbsp;philosophy)
- [Test Suite](Test&nbsp;suite)
## Overview
myGameList is a website built to help gaming enthusiasts to keep track of their gaming backlog and connect with a community while doing so. 
The website is currently self-hosted on a private linux server, and a containerized version is provided for local hosting.

***All game data used in the project is provided by [Giantbomb](https://www.giantbomb.com/)***
## Technologies
**Backend:** Go, testcontainers

**Frontend:** TypeScript, React,  TanStack Router & Query (formerly known as React Query)

**Server:** Docker, Linux, nginx

**CI/CD:** GitLab, GitHub Actions

**Testing:** Vite, Vitest, PlayWright, Go test packages

**Linters:** ESLint/typescript-eslint, Golangci-lint

A full list of dependencies is available in:
- `api/go.mod` for the backend
- `client/package.json` for the frontend
## Requirements
```sh
node 24.0 or newer
go 1.24 or newer
```
In the case of someone hosting their own instance of the app locally, [Docker Desktop](https://docs.docker.com/desktop/setup/install/windows-install/) is recommended.

## Install & Run
**Using Docker**
1. Clone the repository
2. Set up environment variables according to .env.example files
3. Navigate to `/local/` and run:
 `docker compose up -d`
 
4. Application runs `localhost:3004` by default

## Test Suite
Backend tests use a variety of Go packages like mock, assertion, and testcontainers. Integration tests for the backend can be found in `api/integration_test`, unit tests share folders with their UUT  (Unit Under Test) and end in `_test.go`. Not many unit tests were made due to the application being too simple for them to be worth the effort.

Frontend Unit tests are made with Vitest. The tests are located in `client/components`and end in `.test.tsx`. Integration tests were not made due to the simplicity of the application. E2E(End to End) tests were made to test all usage scenarios instead. E2E tests are located in `tests`.

### Install dependencies 
```bash
$ client/ npm install
$ api/ go install
```
### Run frontend unit tests
```bash
$ client/ npm run test
```
### Run frontend E2E tests
```bash
$ npx playwright test (--ui flag for GUI)
```
### Run all Backend tests
```bash
$ api/ go test ./...
```
## Design Philosophy
The main goal of this project is to improve as a software developer professional. There was a need to create something useful and do deployment, pipelines, and architecture from scratch while properly testing all the application features. There was no 3rd party to say what is right and what is wrong but an effort was made in order to keep the codebase maintainable, scalable, and testable.

On the backend, TDD development style was adopted late into the production. External APIs are mocked in tests to account for negative paths. For CI/CD pipelines, GitLab is used because the workflow was compatible with the self hosted deployment solution. Linting is used comprehensively in order to ensure good code quality.
