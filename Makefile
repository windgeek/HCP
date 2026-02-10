.PHONY: all build clean test

all: build

build:
	go build -o hcp ./cmd/hcp

test:
	go test ./pkg/...


clean:
	rm -f hcp
