package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/input"
)

type IGame interface {
	Objects() *ObjectManager
	Keyboard() *input.Keyboard
	Controls() *Controls
}

type Objecter interface {
	Name() string
	ObjType() ObjectType
	Image() *ebiten.Image
	Position() Position
	SetPosition(Position)
	Direction() Direction
	SetDirection(Direction)
	Speed() int
	Update(g IGame, tick uint) error
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

type Breached struct {
	Min struct {
		X, Y bool
	}
	Max struct {
		X, Y bool
	}
}

func (b *Breached) Get() *Breached {
	return b
}
