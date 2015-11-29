# GMon - basic metrics monitoring
# Copyright (c) 2015 Garrett Woodworth (https://github.com/gwoo)

.PHONY: bin/gmon bin/gmon-linux-amd64 bin/gmon-darwin-amd64 test get
VERSION := $(shell git describe --always --dirty --tags)
VERSION_FLAGS := -ldflags "-X main.Version=$(VERSION)"

bin/gmon: bin get
	go build -v $(VERSION_FLAGS) -o $@ .

bin/gmon-linux-amd64: bin get
	GOOS=linux GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ .

bin/gmon-darwin-amd64: bin get
	GOOS=darwin GOARCH=amd64 go build -v $(VERSION_FLAGS) -o $@ .

get:
	go get

test:
	go test -v ./...

bin:
	mkdir bin