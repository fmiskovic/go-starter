build:
	@go build -o bin/app ./cmd/app/

run: build
	@./bin/app

clear:
	@rm -rf bin

test:
	go test -v ./...

race: build
	go test -v ./... --race

cover: 
	go test -cover ./...

css: 
	npm run css

cssi:
	npm install