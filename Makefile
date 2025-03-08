.PHONY: build

build:
	wails build

dev:
	wails dev

test:
	go test ./...

tidy:
	go mod tidy

up:
	go get -u
