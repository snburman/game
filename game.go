package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/input"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000
const MAX_WORKERS = 10

type Game struct {
	tick     uint
	assets   *assets.Assets
	objects  *objects.ObjectManager
	keyboard *input.Keyboard
}

func NewGame() *Game {
	return &Game{
		assets:   assets.Load(),
		objects:  objects.NewObjectManager(),
		keyboard: input.NewKeyboard(),
	}
}

func (g *Game) Update(screen *ebiten.Image) error {
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
		err := object.Update(screen, g.keyboard, g.tick)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Color(color.RGBA{255, 255, 255, 255}))
	objects := g.objects.GetAll()
	for _, o := range objects {
		object := *o
		object.Draw(screen, g.tick)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 336, 336
}

func (g *Game) Run() error {
	return ebiten.RunGame(g)
}

func (g *Game) Assets() *assets.Assets {
	return g.assets
}

func (g *Game) Objects() *objects.ObjectManager {
	return g.objects
}
