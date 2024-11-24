.PHONY: make assets

make:
	go run cmd/magicgame/main.go

build:
	sudo env GOOS=js GOARCH=wasm go build -o game.wasm github.com/snburman/magicgame/cmd/magicgame
	cp game.wasm ../magic_game_client/public/game.wasm

assets:
	go run cmd/generate/main.go

