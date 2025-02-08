package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/snburman/game/config"
)

type Controls struct {
	objs []*Object
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
		Url:  "buttons/a_button.png",
		Name: "aButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 275,
				Y: 365,
			},
			Direction: Right,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "buttons/b_button.png",
		Name: "bButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 225,
				Y: 415,
			},
			Direction: Right,
			Speed:     1,
			Scale:     1,
		},
	},
}

func NewControls() *Controls {
	var objects []*Object
	for _, img := range f {
		objects = append(objects, NewObjectFromFile(img))
	}
	return &Controls{
		objs: objects,
	}
}

func (c *Controls) Objects() []*Object {
	return c.objs
}

var currentIDs = make(map[ebiten.TouchID]bool)

func (c *Controls) Update(g IGame, tick uint) error {
	//TODO: lift logic to input package, touch.go
	// allIDs := []ebiten.TouchID{}
	newPressedIDs := []ebiten.TouchID{}
	justPressedIDs := make(map[ebiten.TouchID]bool)
	justReleasedIDs := make(map[ebiten.TouchID]bool)
	newReleasedIDs := []ebiten.TouchID{}

	newPressedIDs = inpututil.AppendJustPressedTouchIDs(newPressedIDs)
	for _, id := range newPressedIDs {
		justPressedIDs[newPressedIDs[id]] = true
		currentIDs[newPressedIDs[id]] = true
		// allIDs = append(allIDs, id)
	}

	newReleasedIDs = inpututil.AppendJustReleasedTouchIDs(newReleasedIDs)
	for i := range newReleasedIDs {
		justReleasedIDs[newReleasedIDs[i]] = true
		delete(currentIDs, newReleasedIDs[i])
	}

	player := g.Player()
	pos := player.Position()
	// set speed to default
	player.SetSpeed(config.WalkSpeed)
	for _, control := range c.objs {
		for id := range currentIDs {
			x, y := ebiten.TouchPosition(id)
			if control.IsPressed(x, y) {
				switch control.Name() {
				// TODO: create a common interface for this and keyboard
				case "bButton":
					player.SetSpeed(config.RunSpeed)
				case "upButton":
					if !player.Breached().Min.Y {
						player.SetDirection(Up)
						pos.Move(Up, player.Speed())
					}
				case "downButton":
					if !player.Breached().Max.Y {
						player.SetDirection(Down)
						pos.Move(Down, player.Speed())
					}
				case "leftButton":
					if !player.Breached().Min.X {
						player.SetDirection(Left)
						pos.Move(Left, player.Speed())
					}
				case "rightButton":
					if !player.Breached().Max.X {
						player.SetDirection(Right)
						pos.Move(Right, player.Speed())
					}
				}
			}

		}
	}

	return nil
}
