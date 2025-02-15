package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/input"
)

type Player struct {
	*Object
}

func NewPlayer(obj *Object) *Player {
	p := &Player{
		Object: obj,
	}
	p.name = "player"
	return p
}

func (p *Player) Update(g IGame, tick uint) error {
	// check for impending collision with screen boundaries
	p.DetectScreenCollision()
	for _, o := range g.Objects() {
		// if object is a portal, load the map
		if o.ObjType() == ObjectPortal {
			if p.IsCollided(o) && g.CurrentMap().ID.Hex() != o.ID() {
				// if map does not exist, nothing will happen
				g.LoadMap(o.ID())
			}
			continue
		}
		// check for impending collision with other objects
		p.DetectObjectCollision(o)
	}

	var f input.InputFunctions = map[input.Key]func(){
		input.Up: func() {
			p.SetDirection(Up)
			if !p.Breached().Min.Y {
				p.Position().Move(Up, p.Speed())
			}
		},
		input.Down: func() {
			p.SetDirection(Down)
			if !p.Breached().Max.Y {
				p.Position().Move(Down, p.Speed())
			}
		},
		input.Left: func() {
			p.SetDirection(Left)
			if !p.Breached().Min.X {
				p.Position().Move(Left, p.Speed())
			}
		},
		input.Right: func() {
			p.SetDirection(Right)
			if !p.Breached().Max.X {
				p.Position().Move(Right, p.Speed())
			}
		},
	}

	for key, fn := range f {
		if g.Keyboard().IsPressed(key) {
			fn()
		}
	}

	return p.Object.Update(g, tick)
}

func (p *Player) Draw(screen *ebiten.Image, tick uint) {
	p.Object.Draw(screen, tick)
}
