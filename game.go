package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/api"
	"github.com/snburman/game/config"
	"github.com/snburman/game/input"
	"github.com/snburman/game/models"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000

type Game struct {
	debug        *ebiten.Image
	tick         uint
	touchManager *input.TouchManager
	mapService   *api.MapService
	keyboard     *objects.Keyboard
	controls     *objects.Controls
	player       *objects.Player
}

func NewGame() *Game {
	g := &Game{
		debug:        ebiten.NewImage(config.ScreenWidth, 100),
		touchManager: input.NewTouchManager(),
		keyboard:     objects.NewKeyboard(),
		controls:     objects.NewControls(),
	}
	ms := api.NewMapService(api.ApiClient)
	g.mapService = ms
	return g
}

func (g *Game) DebugScreen() *ebiten.Image {
	return g.debug
}

func (g *Game) ClearDebugScreen() {
	g.debug = ebiten.NewImage(config.ScreenWidth, 100)
}

func (g *Game) TouchManager() *input.TouchManager {
	return g.touchManager
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
	for _, p := range g.OnlinePlayers() {
		p.Update(g, g.tick)
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

	// draw game objects
	objects := g.Objects()
	for _, o := range objects {
		o.Draw(screen, g.tick)
	}

	// draw online players
	for _, p := range g.OnlinePlayers() {
		p.Draw(screen, g.tick)
	}

	// draw controls
	for _, o := range g.controls.Objects() {
		o.Draw(screen, g.tick)
	}

	// draw local player
	g.Player().Draw(screen, g.tick)
}

func (g *Game) LoadMap(id string) error {
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

func (g *Game) OnlinePlayers() map[string]*objects.Player {
	return g.mapService.OnlinePlayers()
}

func (g *Game) SetPlayer(player *objects.Player) {
	if player.ObjType() != objects.ObjectPlayer {
		panic("player must have ObjectType: ObjectPlayer")
	}
	g.player = player
}

func (g *Game) DispatchUpdatePlayer() {
	g.mapService.DispatchUpdatePlayer(g.Player())
}

func (g *Game) Keyboard() *objects.Keyboard {
	return g.keyboard
}

func (g *Game) Controls() *objects.Controls {
	return g.controls
}
