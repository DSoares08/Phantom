build:
	go build -buildvcs=false -o ./bin/phantom

run: build
	./bin/phantom

test:
	go test -v ./...
