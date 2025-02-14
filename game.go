package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/api"
	"github.com/snburman/game/config"
	"github.com/snburman/game/models"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000

type Game struct {
	debug      *ebiten.Image
	tick       uint
	mapService *api.MapService
	keyboard   *objects.Keyboard
	controls   *objects.Controls
	player     *objects.Player
}

func NewGame() *Game {
	g := &Game{
		debug:    ebiten.NewImage(config.ScreenWidth, 100),
		keyboard: objects.NewKeyboard(),
		controls: objects.NewControls(),
	}
	ms := api.NewMapService(api.ApiClient)
	g.mapService = ms
	return g
}

func (g *Game) DebugScreen() *ebiten.Image {
	return g.debug
}

func (g *Game) Update() error {
	if g.tick == MAX_TICS {
		g.tick = 1
	} else {
		g.tick++
	}

	objs := g.Objects()
	for _, o := range objs {
		err := o.Update(g, g.tick)
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
	screen.DrawImage(g.DebugScreen(), &ebiten.DrawImageOptions{})
	objects := g.Objects()
	for _, o := range objects {
		o.Draw(screen, g.tick)
	}
	for _, o := range g.controls.Objects() {
		o.Draw(screen, g.tick)
	}
	g.Player().Draw(screen, g.tick)
}

func (g *Game) LoadMap(id string) error {
	// load map
	err := g.mapService.LoadMap(g, id)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) CurrentMap() models.Map[[]models.Image] {
	return g.mapService.CurrentMap()
}

func (g *Game) PrimaryMap() models.Map[[]models.Image] {
	return g.mapService.PrimaryMap()
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

func (g *Game) Objects() []objects.Objecter {
	return g.mapService.CurrentObjects()
}

func (g *Game) Player() *objects.Player {
	return g.mapService.Player()
}

func (g *Game) SetPlayer(player *objects.Player) {
	if player.ObjType() != objects.ObjectPlayer {
		panic("player must have ObjectType: ObjectPlayer")
	}
	g.player = player
}

func (g *Game) Keyboard() *objects.Keyboard {
	return g.keyboard
}

func (g *Game) Controls() *objects.Controls {
	return g.controls
}
