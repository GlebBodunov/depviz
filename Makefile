.PHONY: build run test clean

build:
	go build -o bin/depviz ./src

run:
	go run ./src -config config.yaml

test:
	go test ./tests/ -v

clean:
	rm -rf bin