.PHONY: make assets

make:
	go run cmd/game/main.go

build:
	echo $(date)
	sudo env GOOS=js GOARCH=wasm go build -o game.wasm github.com/snburman/game/cmd/game
	cp game.wasm ../game_server/static/game.wasm