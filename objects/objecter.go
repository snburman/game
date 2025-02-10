package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/config"
)

type IGame interface {
	Objects() *ObjectManager
	LoadMap(id string) error
	CurrentMap() assets.Map[[]assets.Image]
	Player() Objecter
	Keyboard() *Keyboard
	Controls() *Controls
}

type Objecter interface {
	ID() string
	Name() string
	ObjType() ObjectType
	Image() *ebiten.Image
	Position() *Position
	SetPosition(Position)
	Breached() Breached
	Direction() Direction
	SetDirection(Direction)
	Speed() int
	SetSpeed(int)
	Update(g IGame, tick uint) error
	Draw(screen *ebiten.Image, tick uint)
}

// ObjectersFromImages creates a slice of Objecter from a slice of assets.Image
// and returns a pointer to the player object if one exists
func ObjectersFromImages(images []assets.Image) (objs []Objecter, player *Player) {
	for _, img := range images {
		object := NewObject(img, ObjectOptions{
			Position: Position{
				X: img.X,
				Y: img.Y,
			},
			Direction: Right,
			Scale:     config.Scale,
			Speed:     config.WalkSpeed,
		})
		if object.ObjType() == ObjectPlayer {
			if player == nil {
				player = NewPlayer(*object)
			}
			switch img.AssetType {
			case assets.PlayerUp:
				player.SetFrame(Up, object.Image())
				continue
			case assets.PlayerDown:
				player.SetFrame(Down, object.Image())
				continue
			case assets.PlayerLeft:
				player.SetFrame(Left, object.Image())
				continue
			case assets.PlayerRight:
				player.SetFrame(Right, object.Image())
				continue
			}
		}
		objs = append(objs, object)
	}
	return objs, player
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
