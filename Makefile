.PHONY: wasm clean build-mac build-windows build-linux build-all

BUILD_DIR := build
APP_NAME := deadjump

wasm:
	mkdir -p web
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/main.go

build-mac:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 ./cmd/main.go
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-amd64 ./cmd/main.go
	lipo -create -output $(BUILD_DIR)/$(APP_NAME)-darwin-universal $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(BUILD_DIR)/$(APP_NAME)-darwin-amd64
	rm $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 $(BUILD_DIR)/$(APP_NAME)-darwin-amd64

# Linux build requires Docker (Ebiten needs CGO + native libs)
build-linux:
	mkdir -p $(BUILD_DIR)
	docker run --rm \
		-v $(CURDIR):/src \
		-w /src \
		-e CGO_ENABLED=1 \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		golang:1.23 \
		sh -c "apt-get update && apt-get install -y libasound2-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev && go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 ./cmd/main.go"

# Windows build requires Docker with mingw cross-compiler
build-windows:
	mkdir -p $(BUILD_DIR)
	docker run --rm \
		-v $(CURDIR):/src \
		-w /src \
		-e CGO_ENABLED=1 \
		-e GOOS=windows \
		-e GOARCH=amd64 \
		-e CC=x86_64-w64-mingw32-gcc \
		dockcross/windows-static-x64 \
		sh -c "go build -o $(BUILD_DIR)/$(APP_NAME)-windows-amd64.exe ./cmd/main.go"

build-all: build-mac build-windows build-linux

clean:
	rm -rf web/game.wasm $(BUILD_DIR)
