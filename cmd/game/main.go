//go:build js && wasm

package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/objects"
)

func main() {
	game := game.NewGame()

	gameAssets := game.Assets()

	var player *objects.Player
	for _, img := range gameAssets.Images {
		fmt.Println(img.AssetType)
		object := objects.NewObject(img, objects.ObjectOptions{
			Position: objects.Position{
				X: img.X,
				Y: img.Y,
			},
			Direction: objects.Right,
			Scale:     3.5,
			Speed:     3,
		})
		if object.ObjType == objects.ObjectPlayer {
			if player == nil {
				player = objects.NewPlayer(*object)
			}
			switch img.AssetType {
			case assets.PlayerUp:
				player.SetFrame(objects.Up, object.Image())
				continue
			case assets.PlayerDown:
				player.SetFrame(objects.Down, object.Image())
				continue
			case assets.PlayerLeft:
				player.SetFrame(objects.Left, object.Image())
				continue
			case assets.PlayerRight:
				player.SetFrame(objects.Right, object.Image())
				continue
			}
		}
		game.Objects().Add(object)
	}
	if player != nil {
		game.Objects().Add(player)
	} else {
		// TODO: Add default player images if none found
	}

	ebiten.SetScreenTransparent(true)
	ebiten.SetWindowSize(336, 336)
	ebiten.SetWindowTitle("Game")

	if err := game.Run(); err != nil {
		panic(err)
	}
}
