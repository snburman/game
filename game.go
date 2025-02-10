package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/api"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/config"
	"github.com/snburman/game/objects"
)

const MAX_TICS = 10000

type Game struct {
	tick       uint
	assets     *assets.Assets
	objects    *objects.ObjectManager
	mapService *api.MapService
	keyboard   *objects.Keyboard
	controls   *objects.Controls
	player     objects.Objecter
}

func NewGame() *Game {
	g := &Game{
		assets:     &assets.Assets{},
		objects:    objects.NewObjectManager(),
		mapService: api.NewMapService(api.ApiClient),
		keyboard:   objects.NewKeyboard(),
		controls:   objects.NewControls(),
	}
	imgs := g.mapService.ImagesFromMap(g.mapService.PrimaryMap())
	g.assets.Images = imgs
	objs, player := objects.ObjectersFromImages(imgs)
	g.Objects().SetAll(objs)
	if player != nil {
		g.SetPlayer(player)
	}

	// load static objects
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

	// TODO: Update all objects and share updates with server

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
	_map, err := g.mapService.GetMapByID(id)
	if err != nil {
		return err
	}
	// set map
	g.mapService.SetCurrentMap(_map)
	// set images
	images := g.mapService.CurrentImages()
	g.assets.Images = images
	// set objects
	objs, player := objects.ObjectersFromImages(images)
	g.objects.SetAll(objs)
	if player != nil {
		g.SetPlayer(player)
	}
	return nil
}

func (g *Game) CurrentMap() assets.Map[[]assets.Image] {
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
