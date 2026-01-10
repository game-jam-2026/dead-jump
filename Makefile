.PHONY: wasm serve clean

wasm:
	mkdir -p web
	GOOS=js GOARCH=wasm go build -o web/game.wasm ./cmd/main.go

clean:
	rm -rf web/game.wasm
