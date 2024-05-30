# Go compiler
GO := go

# Binary name
BINARY := ./bin/1brc

# Source files
SOURCES := $(wildcard *.go)

.PHONY: all build run test clean

all: build

build: $(BINARY)

$(BINARY): $(SOURCES)
	$(GO) build -o $(BINARY)

run: build
	$(BINARY) -v v1

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY)