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
	cd frontend &&\
		npm install

lint:
	cd frontend &&\
		npm run lint

lintstaged:
	cd frontend &&\
		npm run lint-staged

typecheck:
	cd frontend &&\
		npm run typecheck

up:
	go get -u
