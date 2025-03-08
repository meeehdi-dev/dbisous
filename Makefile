.PHONY: build

build:
	wails build

dev:
	wails dev

test:
	go test ./...

tidy:
	go mod tidy

install:
	cd frontend && npm install

up:
	go get -u
