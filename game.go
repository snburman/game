package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/input"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000

type Game struct {
	tick     uint
	assets   *assets.Assets
	objects  *objects.ObjectManager
	keyboard *input.Keyboard
	controls *objects.Controls
}

func NewGame() *Game {
	return &Game{
		assets:   assets.Load(),
		objects:  objects.NewObjectManager(),
		keyboard: input.NewKeyboard(),
		controls: objects.NewControls(),
	}
}

func (g *Game) Update() error {
	if g.tick > MAX_TICS {
		g.tick = 1
	} else {
		g.tick++
	}

	// TODO: Update all objects and share updates with server

	g.keyboard.Update()
	objects := g.objects.GetAll()
	for _, o := range objects {
		object := *o
		if object.Name() == "player" {

		}
		err := object.Update(g.keyboard, g.tick)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Color(color.White))
	objects := g.objects.GetAll()
	for _, o := range objects {
		object := *o
		object.Draw(screen, g.tick)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 336, 500
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

func (g *Game) Keyboard() *input.Keyboard {
	return g.keyboard
}

func (g *Game) Controls() *objects.Controls {
	return g.controls
}
