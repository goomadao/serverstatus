export PATH := $(PATH):$(GOPATH)/bin
.PHONY: all build test client server web clean

all: test build

build: client server

test:
	go get ./...

client:
	go build -o bin/client ./client/client.go

server: web
	go build -o bin/server ./server/server.go

web:
	cd web; \
	npm install; \
	npm run build

clean:
	rm -f bin/client
	rm -f bin/server
	rm -rf web/dist
	rm -rf web/node_modules