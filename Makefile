all: # this command is like a shortcut that builds the app, runs db migration, and runs the app server
	@$(MAKE) build
	@$(MAKE) db cmd=init
	@$(MAKE) db cmd=migrate
	@$(MAKE) run

build: # build the go code
	@echo "building app started"
	@go build -o bin/app ./cmd/app/
	@echo "building app finished"

run: # run server
	@echo "starting the app..."
	@./bin/app serve

db: # db migration related commands, like init, migrate, status, rollback...
	@./bin/app db $(cmd)

run-db: # run posgres db in docker with default configs
	@echo "starting go-db..."
	@docker run --name go-db -e POSTGRES_PASSWORD=dbadmin -e POSTGRES_USER=dbadmin -e PGDATA=/var/lib/postgresql/data -e POSTGRES_DB=go-db --volume=/var/lib/postgresql/data -p 5432:5432 -d postgres
	@echo "posgres go-db started."

clear: # delete app build
	@rm -rf bin

test: # run tests
	@go test -v ./...

race: # check race conditions
	@go test -v ./... --race

cover: # check test coverage
	@go test -cover ./...

css: # download taiwind css in ./public/assets/app.css
	@echo "running npm css..."
	npm run css

cssi: # install tailwind css
	@echo "running npm install..."
	@npm install

go-format: # format go code
	@go fmt ./...

go-lint: # check go lint
	@echo "golangci-lint run..."
	@golangci-lint run --timeout 5m
	@echo "done"

