# Go Starter
A fullstack starter pack for default ssr rendering, inspired by [anthdm's boredstack](https://github.com/anthdm/boredstack).

The full-stack is: GO, Fiber, Bun, Postgres, HTML, Tailwind CSS and AlpineJS.

### Style
- `make cssi` - runs `npm install` and install tailwind
- `make css` - generates css file: ./public/assets/app.css


### Build and Run commands
- `make run` - will build and run everything, visit address: `http://localhost:8080`

### Environment variables
- `HTTP_LISTEN_ADDR`  - default is `:8080`
- `PRODUCTION` - default is `false`
- `DB_PASSWORD`
- `DB_USER`
- `DB_NAME`
- `DB_HOST`

### TODO
- Add integration db tests
- Add integration handler tests
- Swagger for API
- Add AlpineJS entity handling
