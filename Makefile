.PHONY: build-app run-app

APP_LINUX_BIN=build/app-linux-amd64

build-app:
	GOOS=linux GOARCH=amd64 go build -o $(APP_LINUX_BIN) cmd/app/main.go

run-app:
	go run ./cmd/app/main.go -c ./config.local.yml
