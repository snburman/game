package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/config"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000

type Game struct {
	tick     uint
	assets   *assets.Assets
	objects  *objects.ObjectManager
	keyboard *objects.Keyboard
	controls *objects.Controls
	player   objects.Objecter
}

func NewGame() *Game {
	return &Game{
		assets:   assets.Load(),
		objects:  objects.NewObjectManager(),
		keyboard: objects.NewKeyboard(),
		controls: objects.NewControls(),
	}
}

func (g *Game) Update() error {
	if g.tick == MAX_TICS {
		g.tick = 1
	} else {
		g.tick++
	}

	// TODO: Update all objects and share updates with server

	objs := g.objects.GetAll()
	for _, o := range objs {
		object := *o
		if object.ObjType() == objects.ObjectPlayer {

		}
		err := object.Update(g, g.tick)
		if err != nil {
			return err
		}
	}
	g.controls.Update(g, g.tick)
	g.keyboard.Update(g)
	g.Player().Update(g, g.tick)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Color(color.RGBA{
		175, 175, 178, 255,
	}))
	objects := g.objects.GetAll()
	for _, o := range objects {
		object := *o
		object.Draw(screen, g.tick)
	}
	for _, o := range g.controls.Objects() {
		object := *o
		object.Draw(screen, g.tick)
	}
	g.Player().Draw(screen, g.tick)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}

func (g *Game) Run() error {
	return ebiten.RunGame(g)
}

func (g *Game) RunGameWithOptions(opts *ebiten.RunGameOptions) error {
	return ebiten.RunGameWithOptions(g, opts)
}

func (g *Game) Assets() *assets.Assets {
	return g.assets
}

func (g *Game) Objects() *objects.ObjectManager {
	return g.objects
}

func (g *Game) Player() objects.Objecter {
	return g.player
}

func (g *Game) SetPlayer(player objects.Objecter) {
	g.player = player
}

func (g *Game) Keyboard() *objects.Keyboard {
	return g.keyboard
}

func (g *Game) Controls() *objects.Controls {
	return g.controls
}
