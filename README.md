![go workflow](https://github.com/fmiskovic/go-starter/actions/workflows/go-ci.yml/badge.svg)
![lint workflow](https://github.com/fmiskovic/go-starter/actions/workflows/golangci-lint.yml/badge.svg)
___
### ⚠️ WORK IN PROGRESS ⚠️
___
# Golang Fullstack Starter Pack

## Description
A comprehensive starter pack for Golang fullstack development, designed to streamline the process of setting up a fullstack web application with Server-Side Rendering (SSR) out of the box.

## Motivation
As a backend developer, I sought a suitable solution for the frontend and experimented with React, Vue (both with and without Nuxt). Despite their unnecessary complexity, numerous dependencies for a simple "Hello world," and the overhead introduced by features like the virtual DOM, I also observed a trend toward server-side rendering (SSR).

In light of this shift, let me express that we don't necessarily require React, Vue, Svelte, Angular, or any similar frameworks. What we truly need is simplicity.

## Why did I build this?
As a developer, I often found myself spending valuable time on repetitive tasks when starting a new Golang fullstack project. This starter pack is my solution to this problem, offering a pre-configured environment that allows developers to focus on building features rather than dealing with boilerplate code and setup complexities.

## What problem does it solve?
This project addresses the pain points associated with the initial setup of fullstack applications. It provides a ready-to-use template with SSR support, allowing developers to kickstart their projects without the hassle of configuring the frontend and backend separately.

## How to use this template
> **DO NOT FORK** this is meant to be used from **[Use this template](https://github.com/fmiskovic/go-starter/generate)** feature.

## How to build and run the project
Since it is using postgres db, pre-condition is to have running postgres.
If you want to run the app for the first time, follow the steps bellow, and it will bi accessible in [http://localhost:8080](http://localhost:8080)

1) Run postgres db: ```docker run --name go-db -e POSTGRES_PASSWORD=dbadmin -e POSTGRES_USER=dbadmin -e PGDATA=/var/lib/postgresql/data -e POSTGRES_DB=go-db --volume=/var/lib/postgresql/data -p 5432:5432 -d postgres```
2) Build the app: ```make build```
3) Init db migration: ```make db cmd=init```
4) Migrate db tables: ```make db cmd=migrate```
5) Run the app: ```make run```

## Available commands and variables

### Database migration commands
- `make db cmd=init` - initialize bun migration
- `make db cmd=migrate` - applies migration scripts (files with .sql extension)
- `make db cmd=status` - check migration status
- `make db cmd=rollback` - rollback last migration group

### CSS commands
- `make cssi` - runs `npm install` and install tailwind
- `make css` - generates css file: ./public/assets/app.css

### Variables
- `HTTP_LISTEN_ADDR`  - default is `:8080`
- `PRODUCTION` - default is `false`
- `DB_PASSWORD` - default is `dbadmin`
- `DB_USER` - default is `dbadmin`
- `DB_NAME` - default is `go-db`
- `DB_HOST` - defailt is `localhost`

### TODO list
- Unit tests
- Integration tests for handlers
- Basic authentication and authorization
- Swagger for API
- Add AlpineJS entity handling (CRUD operations)
- Users view

## License
This project is licensed under the [MIT License](https://github.com/fmiskovic/go-starter/blob/main/LICENSE.md).

### Credits
Inspired by [anthdm's boredstack](https://github.com/anthdm/boredstack).


