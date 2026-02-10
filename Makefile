.PHONY: all build clean test

all: build


build:
	go build -o hcp ./cmd/hcp
	go build -o hcp-release ./cmd/hcp-release

test:
	go test ./pkg/...


clean:
	rm -f hcp
