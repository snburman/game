package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/input"
)

type Player struct {
	Object
}

func NewPlayer(obj Object) *Player {
	p := &Player{
		Object: obj,
	}
	p.name = "player"
	return p
}

func (p *Player) Update(g IGame, tick uint) error {
	p.DetectScreenCollision()
	for _, o := range g.Objects().GetAll() {
		p.DetectObjectCollision(*o)
	}

	pos := p.Position()
	var f input.InputFunctions = map[input.Key]func(){
		input.Up: func() {
			if !p.Breached().Min.Y {
				pos.Move(Up, p.Speed())
				p.SetDirection(Up)
			}
		},
		input.Down: func() {
			if !p.Breached().Max.Y {
				pos.Move(Down, p.Speed())
				p.SetDirection(Down)
			}
		},
		input.Left: func() {
			if !p.Breached().Min.X {
				pos.Move(Left, p.Speed())
				p.SetDirection(Left)
			}
		},
		input.Right: func() {
			if !p.Breached().Max.X {
				pos.Move(Right, p.Speed())
				p.SetDirection(Right)
			}
		},
	}

	for key, fn := range f {
		if g.Keyboard().IsPressed(key) {
			fn()
		}
	}

	p.SetPosition(*pos)

	return p.Object.Update(g, tick)
}

func (p *Player) Draw(screen *ebiten.Image, tick uint) {
	p.Object.Draw(screen, tick)
}
