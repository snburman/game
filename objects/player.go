package objects

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game/input"
)

// Image frame indices
const (
	FaceUp = iota
	FaceDown
	FaceLeft
	FaceRight
)

type Player struct {
	Object
}

func NewPlayer(obj Object) *Player {
	return &Player{
		Object: obj,
	}
}

func (p *Player) Update(screen *ebiten.Image, i input.Input, tick uint) error {
	pos := p.Position()

	var f input.InputFunctions = map[input.Key]func(){
		input.Up: func() {
			pos.Move(Up, p.Speed())
			p.SetDirection(Up)
			p.SetCurrentFrame(FaceUp)
		},
		input.Down: func() {
			pos.Move(Down, p.Speed())
			p.SetDirection(Down)
			p.SetCurrentFrame(FaceDown)
		},
		input.Left: func() {
			pos.Move(Left, p.Speed())
			p.SetDirection(Left)
			p.SetCurrentFrame(FaceLeft)
		},
		input.Right: func() {
			pos.Move(Right, p.Speed())
			p.SetDirection(Right)
			p.SetCurrentFrame(FaceRight)
		},
	}

	for key, fn := range f {
		if i.IsPressed(key) {
			fn()
		}
	}

	p.SetPosition(pos)

	return p.Object.Update(screen, i, tick)
}

func (p *Player) Draw(screen *ebiten.Image, tick uint) {
	p.Object.Draw(screen, tick)
}
