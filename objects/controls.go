package objects

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Controls struct {
	objs []*Object
}

var TouchIDs = []ebiten.TouchID{}
var Touches = make(map[ebiten.TouchID]*Touch, 0)

type Touch struct {
	OriginX, OriginY int
	CurrX, CurrY     int
	Duration         int
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

func (c *Controls) Update(g IGame, tick uint) error {
	for id := range Touches {
		if inpututil.IsTouchJustReleased(id) {
			fmt.Println("Touch released: " + fmt.Sprint(id))
			delete(Touches, id)
		}
	}

	fmt.Println(TouchIDs)
	TouchIDs = inpututil.AppendJustPressedTouchIDs(TouchIDs)
	player := g.Player()
	for _, id := range TouchIDs {
		x, y := ebiten.TouchPosition(id)
		Touches[id] = &Touch{
			OriginX: x, OriginY: y,
			CurrX: x, CurrY: y,
		}
		pos := player.Position()
		for _, control := range c.objs {
			if control.IsPressed(id) {
				switch control.Direction() {
				case Up:
					fmt.Println("UP")
					if !player.Breached().Min.Y {
						player.SetDirection(Up)
						pos.Move(Up, player.Speed())
					}
				case Down:
					fmt.Println("DOWN")
					if !player.Breached().Max.Y {
						player.SetDirection(Down)
						pos.Move(Down, player.Speed())
					}
				case Left:
					fmt.Println("LEFT")
					if !player.Breached().Min.X {
						player.SetDirection(Left)
						pos.Move(Left, player.Speed())
					}
				case Right:
					fmt.Println("RIGHT")
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
