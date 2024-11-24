package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/magicgame"
	"github.com/snburman/magicgame/objects"
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
	// exportJSFunctions()
	game := magicgame.NewGame()
	defer game.Stop()

	assets := game.Assets()
	img := assets.Sprite("untitled")
	object := objects.NewObject(img, objects.ObjectOptions{
		Position: objects.Position{
			X: 100,
			Y: 100,
		},
		Direction: objects.Right,
		Scale:     3,
		Speed:     3,
	})
	player := objects.NewPlayer(*object)
	game.Objects().Add(player)

	// amethyst := assets.Sprite("amethyst")

	// go func() {
	// 	// after 5 seconds, add a new object
	// 	<-time.After(5 * time.Second)
	// 	obj2 := objects.NewObject(amethyst, objects.ObjectOptions{
	// 		Position: objects.Position{
	// 			X: 200,
	// 			Y: 200,
	// 		},
	// 		Scale: 2,
	// 	})
	// 	game.Objects().Add(obj2)
	// }()

	ebiten.SetWindowSize(528, 528)
	ebiten.SetWindowTitle("Magic Game")
	ebiten.SetWindowResizable(true)

	if err := game.Run(); err != nil {
		panic(err)
	}
}
