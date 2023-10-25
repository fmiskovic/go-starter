build:
	@go build -o bin/app ./cmd/

run: build
	@./bin/app serve -addr=$(addr)

db: build
	@./bin/app db $(cmd)

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