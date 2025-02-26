package objects

import (
	"github.com/snburman/game/config"
)

type Controls struct {
	objs []*Object
}

type Touch struct {
	X, Y int
}

var buttonImgs []FileImage = []FileImage{
	{
		Url:  "buttons/up_button.png",
		Name: "upButton",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 60,
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
				X: 60,
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
				X: 22,
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
				X: 97,
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
				X: 167,
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
				X: 285,
				Y: 362,
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
				X: 235,
				Y: 418,
			},
			Direction: Right,
			Speed:     1,
			Scale:     1,
		},
	},
}

func NewControls() *Controls {
	var objects []*Object
	for _, img := range buttonImgs {
		objects = append(objects, NewObjectFromFile(img))
	}
	return &Controls{
		objs: objects,
	}
}

func (c *Controls) Objects() []*Object {
	return c.objs
}

func (c *Controls) Update(g IGame, tick uint) error {
	for _, control := range c.objs {
		if control.IsClicked() {
			c.onPress(g, control)
			continue
		}
		for _, touch := range g.TouchManager().Touches() {
			if control.name == "bButton" {
				if control.IsPressed(touch.X, touch.Y) {
					g.Player().SetSpeed(config.RunSpeed)
				} else {
					g.Player().SetSpeed(config.WalkSpeed)
				}
			}
			if control.IsPressed(touch.X, touch.Y) {
				c.onPress(g, control)
			}
		}
	}
	return nil
}

func (c *Controls) onPress(g IGame, control *Object) {
	speed := g.Player().Speed()
	switch control.name {
	case "home_button":
		g.LoadMap(g.PrimaryMap().ID.Hex())
	case "upButton":
		g.Player().SetDirection(Up)
		if !g.Player().Breached().Min.Y {
			g.Player().Position().Move(Up, speed)
		}
		g.DispatchUpdatePlayer()
	case "downButton":
		g.Player().SetDirection(Down)
		if !g.Player().Breached().Max.Y {
			g.Player().Position().Move(Down, speed)
		}
		g.DispatchUpdatePlayer()
	case "leftButton":
		g.Player().SetDirection(Left)
		if !g.Player().Breached().Min.X {
			g.Player().Position().Move(Left, speed)
		}
		g.DispatchUpdatePlayer()
	case "rightButton":
		g.Player().SetDirection(Right)
		if !g.Player().Breached().Max.X {
			g.Player().Position().Move(Right, speed)
		}
		g.DispatchUpdatePlayer()
	default:
	}
}
