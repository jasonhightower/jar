GO ?= go

all: build

install: test
	$(GO) install ./...

test:
	$(GO) test ./...

build:
	cd cmd; $(GO) build
