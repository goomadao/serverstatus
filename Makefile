export PATH := $(PATH):$(GOPATH)/bin
.PHONY: all build test client server web clean

all: test build

build: client server

test:
	go get ./...

client:
	go build -o bin/client ./client/client.go

server: 
	go build -o bin/server ./server/server.go

web:
	cd web; \
	npm install; \
	npm run build

file: web
	go get github.com/rakyll/statik
	rm -rf assets/static
	cp -rf web/dist/static assets
	rm -rf assets/statik
	go generate ./assets/...

clean:
	rm -rf bin
	rm -rf web/dist
	rm -rf web/node_modules