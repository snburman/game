//go:build js && wasm

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/config"
	"github.com/snburman/game/objects"
)

func main() {
	game := game.NewGame()

	gameAssets := game.Assets()

	var player *objects.Player
	for _, img := range gameAssets.Images {
		object := objects.NewObject(img, objects.ObjectOptions{
			Position: objects.Position{
				X: img.X,
				Y: img.Y,
			},
			Direction: objects.Right,
			Scale:     config.Scale,
			Speed:     3,
		})
		if object.ObjType() == objects.ObjectPlayer {
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

	// load static objects
	for _, f := range objects.StaticImages {
		o := objects.NewObjectFromFile(f)
		game.Objects().Add(o)
	}

	// load controls
	controls := objects.NewControls()
	for _, o := range controls.Objects() {
		game.Objects().Add(o)
	}

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
	}
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Game")

	if err := game.RunGameWithOptions(opts); err != nil {
		panic(err)
	}
}
