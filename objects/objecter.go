package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/input"
)

type Objecter interface {
	Name() string
	Position() Position
	SetPosition(Position)
	Direction() Direction
	SetDirection(Direction)
	Speed() int
	Update(input input.Input, tick uint) error
	Draw(screen *ebiten.Image, tick uint)
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Position struct {
	X, Y, Z int
}

func (p *Position) Move(d Direction, s int) {
	switch d {
	case Up:
		p.Y -= s
	case Down:
		p.Y += s
	case Left:
		p.X -= s
	case Right:
		p.X += s
	}
}

func (p *Position) Set(x, y, z int) {
	p.X = x
	p.Y = y
	p.Z = z
}

func WillCollide(source, destination Object) (top, bottom, left, right bool) {
	// TODO: Implement collision detection
	return
}
