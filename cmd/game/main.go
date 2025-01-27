package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game"
	"github.com/snburman/game/objects"
)

func main() {
	// TODO: use global.JS to get user ID
	game := game.NewGame()

	assets := game.Assets()
	for _, img := range assets.Images {
		object := objects.NewObject(img, objects.ObjectOptions{
			Position: objects.Position{
				X: img.X,
				Y: img.Y,
			},
			Direction: objects.Right,
			Scale:     3.5,
			Speed:     3,
		})
		game.Objects().Add(object)
	}
	// player := objects.NewPlayer(*object)

	ebiten.SetWindowSize(336, 336)
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowResizable(true)

	if err := game.Run(); err != nil {
		panic(err)
	}
}
