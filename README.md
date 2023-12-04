![go workflow](https://github.com/fmiskovic/go-starter/actions/workflows/go-ci.yml/badge.svg)
![lint workflow](https://github.com/fmiskovic/go-starter/actions/workflows/golangci-lint.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/fmiskovic/go-starter)](https://goreportcard.com/report/github.com/fmiskovic/go-starter)
___
# Golang Fullstack Starter Pack

## Description
A comprehensive starter pack for Golang fullstack development, hexagonal architecture designed to streamline the process of setting up a fullstack web application with Server-Side Rendering (SSR) and JWT authentication out of the box.

## Motivation
As a backend developer, I sought a suitable solution for the frontend and experimented with React, Vue (with and without Nuxt). Despite their unnecessary complexity, numerous dependencies for a simple "Hello world," and the overhead introduced by features like the virtual DOM, I also observed a trend toward server-side rendering (SSR).

In light of this shift, let me express that we don't necessarily require React, Vue, Svelte, Angular, or any similar frameworks. What we truly need is simplicity.

## Why did I build this?
I often found myself spending valuable time on repetitive tasks when starting a new Golang fullstack project. This starter pack is my solution to this problem, offering a pre-configured environment that allows developers to focus on building features rather than dealing with boilerplate code and setup complexities.

## What problem does it solve?
This project addresses the pain points associated with the initial setup of fullstack applications. It provides a ready-to-use template with hexagonal architecture, SSR and JWT auth support, allowing developers to kickstart their projects without the hassle of configuring the frontend and backend separately.

## Tech Stack
- Backend:
    - [Go](https://go.dev/), [Fiber](https://gofiber.io/), [Bun](https://bun.uptrace.dev/), [Postgres](https://www.postgresql.org/) 
- Frontend:
    - [HTML](https://developer.mozilla.org/en-US/docs/Web/HTML), [Tailwind CSS](https://flowbite.com/), [AlpineJS](https://alpinejs.dev/)

## How to use this template
> **DO NOT FORK** this is meant to be used from **[Use this template](https://github.com/fmiskovic/go-starter/generate)** feature.

## How to build and run the project
Since it is using postgres db, pre-condition is to have running postgres.
If you want to run the app for the first time, follow the steps bellow, and it will bi accessible here: [http://localhost:8080](http://localhost:8080) and Swagger docs here: [http://localhost:8080/api/v1/docs](http://localhost:8080/api/v1/docs)

1) Run postgres db:```make run-db``` OR ```docker run --name go-db -e POSTGRES_PASSWORD=dbadmin -e POSTGRES_USER=dbadmin -e PGDATA=/var/lib/postgresql/data -e POSTGRES_DB=go-db --volume=/var/lib/postgresql/data -p 5432:5432 -d postgres```
2) Build the app, migrate the db, and run the server: ```make all```

## Available commands and variables

### Database migration commands
- `make db cmd=init` - initialize bun migration
- `make db cmd=migrate` - applies migration scripts (files with .sql extension)
- `make db cmd=status` - check migration status
- `make db cmd=rollback` - rollback last migration group

### CSS commands
- `make cssi` - runs `npm install` and install tailwind
- `make css` - generates css file: ./public/assets/app.css

### Other available commands
Look at [Makefile](https://github.com/fmiskovic/go-starter/blob/main/Makefile)

### Variables
- `HTTP_LISTEN_ADDR`  - default is ***:8080***
- `PRODUCTION` - default is ***false***
- `DB_PASSWORD` - default is ***dbadmin***
- `DB_USER` - default is ***dbadmin***
- `DB_NAME` - default is ***go-db***
- `DB_HOST` - defailt is ***localhost:5432***
- `DB_MAX_IDLE_CONN` - default is ***num of cpu + 1***
- `DB_MAX_OPEN_CONN` - default is ***num of cpu + 1***
- `AUTH_JWT_EXP_TIME` - default is ***24 hours***
- `AUTH_JWT_SECRET` - default is ***secret***

### TODO list
- Add AlpineJS entity handling (CRUD operations)
- Users view
- Email notifications
- Push notifications
- i18n

## Contributing

Pull requests are welcome. For major changes, please [open an issue](https://github.com/fmiskovic/go-starter/issues/new) first to discuss what you would like to change.

Please make sure that tests and lint checks are passing, and that your changes are well-tested.

Thank you for contributing!

## License
This project is licensed under the [MIT License](https://github.com/fmiskovic/go-starter/blob/main/LICENSE.md).

## Credits
Inspired by [anthdm's boredstack](https://github.com/anthdm/boredstack).


