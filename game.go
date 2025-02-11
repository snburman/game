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
	tick       uint
	objects    *objects.ObjectManager
	mapService *api.MapService
	keyboard   *objects.Keyboard
	controls   *objects.Controls
	player     objects.Objecter
}

func NewGame() *Game {
	g := &Game{
		objects:  objects.NewObjectManager(),
		keyboard: objects.NewKeyboard(),
		controls: objects.NewControls(),
	}
	ms := api.NewMapService(api.ApiClient, g)
	g.mapService = ms

	// load static images/objects
	for _, f := range objects.StaticImages {
		o := objects.NewObjectFromFile(f)
		g.Objects().Add(o)
	}
	return g
}

func (g *Game) Update() error {
	if g.tick == MAX_TICS {
		g.tick = 1
	} else {
		g.tick++
	}

	objs := g.objects.GetAll()
	for _, o := range objs {
		object := o
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
		object := o
		object.Draw(screen, g.tick)
	}
	for _, o := range g.controls.Objects() {
		object := o
		object.Draw(screen, g.tick)
	}
	g.Player().Draw(screen, g.tick)
}

func (g *Game) LoadMap(id string) error {
	// fetch map by id
	_map, err := g.mapService.GetMapByID(id)
	if err != nil {
		return err
	}
	// set map
	g.mapService.SetCurrentMap(g, _map)
	return nil
}

func (g *Game) CurrentMap() models.Map[[]models.Image] {
	return g.mapService.CurrentMap()
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

func (g *Game) Objects() *objects.ObjectManager {
	return g.objects
}

func (g *Game) Player() *objects.Player {
	player, ok := g.player.(*objects.Player)
	if !ok {
		panic("player must be a pointer to Player")
	}
	return player
}

func (g *Game) SetPlayer(player objects.Objecter) {
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
