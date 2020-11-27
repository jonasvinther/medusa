.PHONY: test build compile start-vault stop-vault

test:
	go test -v ./...

build:
	go build

compile:
	echo "Compiling for every OS and Platform"
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/medusa-linux-amd64 main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o bin/medusa-linux-arm main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/medusa-linux-arm64 main.go
	CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -o bin/medusa-freebsd-386 main.go
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o bin/medusa.exe main.go

start-vault:
	./scripts/start-vault.sh

stop-vault:
	./scripts/stop-vault.sh

all: test build