GO ?= go

all: build

install: test
	$(GO) install ./...

test:
	$(GO) test ./...

build:
	$(GO) build
