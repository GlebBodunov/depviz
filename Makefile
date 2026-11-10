.PHONY: build run test examples clean

build:
	go build -o bin/depviz ./src

run:
	go run ./src -config config.yaml

test:
	go test ./tests/ -v

examples:
	for f in examples/*.yaml; do go run ./src -config $$f || exit 1; done

clean:
	rm -rf bin