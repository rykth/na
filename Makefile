BINARY := bin/na
CMD    := ./cmd/na

.PHONY: build run test lint clean

build:
	go build -o $(BINARY) $(CMD)

run:
	go run $(CMD)

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
