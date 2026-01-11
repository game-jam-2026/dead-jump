.PHONY: wasm clean build-mac build-windows build-linux build-all

BUILD_DIR := build
APP_NAME := deadjump

wasm:
	mkdir -p web
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/main.go

clean:
	rm -rf web/game.wasm $(BUILD_DIR)

build-windows:
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME).exe ./cmd/main.go
