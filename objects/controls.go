package objects

import "github.com/snburman/game/input"

type Controls struct {
	objs []*Object
}

var f []FileImage = []FileImage{
	{
		Url:  "square.png",
		Name: "dPadUp",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 63,
				Y: 360,
			},
			Direction: Up,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "square.png",
		Name: "dPadDown",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 63,
				Y: 440,
			},
			Direction: Down,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "square.png",
		Name: "dPadLeft",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 25,
				Y: 400,
			},
			Direction: Left,
			Speed:     1,
			Scale:     1,
		},
	},
	{
		Url:  "square.png",
		Name: "dPadRight",
		Opts: ObjectOptions{
			ObjectType: ObjectTile,
			Position: Position{
				X: 100,
				Y: 400,
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

func (c *Controls) Update(i input.Input, tick uint) error {
	return nil
}
