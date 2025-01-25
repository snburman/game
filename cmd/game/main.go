package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game"
	"github.com/snburman/game/objects"
)

type Level struct {
	Width   int
	Height  int
	objects []*objects.Object
}

func NewLevel(width, height int, objects []*objects.Object) *Level {
	return &Level{
		Width:   width,
		Height:  height,
		objects: objects,
	}
}

func main() {
	// TODO: use global.JS to get user ID
	game := game.NewGame()

	assets := game.Assets()
	for _, img := range assets.Images.Sprites {
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
	fmt.Println(game.Objects())
	// player := objects.NewPlayer(*object)

	ebiten.SetWindowSize(336, 336)
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowResizable(true)

	if err := game.Run(); err != nil {
		panic(err)
	}
}
