//go:build js && wasm

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game"
	"github.com/snburman/game/config"
)

func main() {
	game := game.NewGame()
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Game")

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}
	if err := game.RunGameWithOptions(opts); err != nil {
		panic(err)
	}
}
