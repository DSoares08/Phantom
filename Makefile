build:
	go build -o ./bin/phantom

run: build
	./bin/phantom

test:
	go test -v ./...
