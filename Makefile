.PHONY: build test tidy clean

BINARY ?= nginx-autoindex-tui

build:
	go build -ldflags "-w -s" -o $(BINARY) .

test:
	go test ./tests/... -v

tidy:
	go mod tidy

clean:
	rm -f $(BINARY)
	go clean -cache
