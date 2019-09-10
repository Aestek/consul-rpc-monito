.ALL: build
.PHONY: build

build:
	go build -ldflags '-w -extldflags "-static"'