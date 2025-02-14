package objects

import (
	"strconv"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/snburman/game/config"
)

type Controls struct {
	mu         sync.Mutex
	currentIDs map[ebiten.TouchID]bool
	objs       []*Object
}

type Touch struct {
	OriginX, OriginY int
	CurrX, CurrY     int
	Duration         int
}

var f []FileImage = []FileImage{
	{
		Url:  "buttons/up_button.png",
		Name: "upButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 63,
				Y: 350,
			},
			Direction: Up,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "buttons/down_button.png",
		Name: "downButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 63,
				Y: 430,
			},
			Direction: Down,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "buttons/left_button.png",
		Name: "leftButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 25,
				Y: 390,
			},
			Direction: Left,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "buttons/right_button.png",
		Name: "rightButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 100,
				Y: 390,
			},
			Direction: Right,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "buttons/home_button.png",
		Name: "home_button",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 250,
				Y: 390,
			},
			Direction: Right,
			Speed:     1,
			Scale:     1,
		},
	},
	// {
	// 	Url:  "buttons/a_button.png",
	// 	Name: "aButton",
	// 	Opts: ObjectOptions{
	// 		ObjectType: ObjectTile,
	// 		Position: Position{
	// 			X: 275,
	// 			Y: 365,
	// 		},
	// 		Direction: Right,
	// 		Speed:     1,
	// 		Scale:     1,
	// 	},
	// },
	// {
	// 	Url:  "buttons/b_button.png",
	// 	Name: "bButton",
	// 	Opts: ObjectOptions{
	// 		ObjectType: ObjectTile,
	// 		Position: Position{
	// 			X: 225,
	// 			Y: 415,
	// 		},
	// 		Direction: Right,
	// 		Speed:     1,
	// 		Scale:     1,
	// 	},
	// },
}

func NewControls() *Controls {
	var objects []*Object
	for _, img := range f {
		objects = append(objects, NewObjectFromFile(img))
	}
	return &Controls{
		objs:       objects,
		currentIDs: make(map[ebiten.TouchID]bool),
	}
}

func (c *Controls) Objects() []*Object {
	return c.objs
}

func (c *Controls) Update(g IGame, tick uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	str := ""
	for id := range c.currentIDs {
		str += strconv.Itoa(int(id)) + "\n"
	}
	ebitenutil.DebugPrint(g.DebugScreen(), str)

	//TODO: lift logic to input package, touch.go
	// allIDs := []ebiten.TouchID{}
	newPressedIDs := []ebiten.TouchID{}
	justPressedIDs := make(map[ebiten.TouchID]bool)
	justReleasedIDs := make(map[ebiten.TouchID]bool)
	newReleasedIDs := []ebiten.TouchID{}

	newPressedIDs = inpututil.AppendJustPressedTouchIDs(newPressedIDs)
	for _, id := range newPressedIDs {
		justPressedIDs[newPressedIDs[id]] = true
		c.currentIDs[newPressedIDs[id]] = true
		// allIDs = append(allIDs, id)
	}

	newReleasedIDs = inpututil.AppendJustReleasedTouchIDs(newReleasedIDs)
	for i := range newReleasedIDs {
		justReleasedIDs[newReleasedIDs[i]] = true
		delete(c.currentIDs, newReleasedIDs[i])
	}

	speed := g.Player().Speed()
	// set speed to default
	g.Player().SetSpeed(config.WalkSpeed)
	for _, control := range c.objs {
		for id := range c.currentIDs {
			x, y := ebiten.TouchPosition(id)
			if control.IsPressed(x, y) {
				switch control.Name() {
				case "home_button":
					g.LoadMap(g.PrimaryMap().ID.Hex())
				case "upButton":
					g.Player().SetDirection(Up)
					if !g.Player().Breached().Min.Y {
						g.Player().Position().Move(Up, speed)
					}
				case "downButton":
					g.Player().SetDirection(Down)
					if !g.Player().Breached().Max.Y {
						g.Player().Position().Move(Down, speed)
					}
				case "leftButton":
					g.Player().SetDirection(Left)
					if !g.Player().Breached().Min.X {
						g.Player().Position().Move(Left, speed)
					}
				case "rightButton":
					g.Player().SetDirection(Right)
					if !g.Player().Breached().Max.X {
						g.Player().Position().Move(Right, speed)
					}
				}
			}

		}
	}

	return nil
}
