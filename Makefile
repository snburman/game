.PHONY: make assets

make:
	go run cmd/game/main.go

build:
	sudo env GOOS=js GOARCH=wasm go build -o game.wasm github.com/snburman/game/cmd/game
	cp game.wasm ../magic_game_client/public/game.wasm

assets:
	go run cmd/generate/main.go

