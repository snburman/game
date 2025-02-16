package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/config"
	"github.com/snburman/game/input"
	"github.com/snburman/game/models"
)

type IGame interface {
	DebugScreen() *ebiten.Image
	ClearDebugScreen()
	TouchManager() *input.TouchManager
	PrimaryMap() models.Map[[]models.Image]
	CurrentMap() models.Map[[]models.Image]
	LoadMap(id string) error
	Objects() []Objecter
	Player() *Player
	SetPlayer(*Player)
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

// ObjectersFromImages creates a slice of Objecter from a slice of models.Image
// and returns a pointer to the player object if one exists
func ObjectersFromImages(images []models.Image, userID string) (objs []Objecter, player *Player) {
	var p *Player
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
			if img.UserID == userID {
				if p == nil {
					p = NewPlayer(object, userID)
				}
				dir := DirectionFromAssetType(img.AssetType)
				p.SetFrame(dir, object.Image())
				continue
			}
			continue
		}
		objs = append(objs, object)
	}
	return objs, p
}

func PlayersFromImages(images []models.Image) map[string]*Player {
	players := make(map[string]*Player)
	for _, img := range images {
		if img.AssetType != models.PlayerUp && img.AssetType != models.PlayerDown &&
			img.AssetType != models.PlayerLeft && img.AssetType != models.PlayerRight {
			continue
		}
		object := NewObject(img, ObjectOptions{
			Position: Position{
				X: img.X,
				Y: img.Y,
			},
			Direction: Right,
			Scale:     config.Scale,
			Speed:     config.WalkSpeed,
		})
		if _, ok := players[img.UserID]; !ok {
			players[img.UserID] = NewPlayer(object, img.UserID)
		}
		players[img.UserID].SetFrame(DirectionFromAssetType(img.AssetType), object.Image())
	}
	return players
}

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func DirectionFromAssetType(t models.AssetType) Direction {
	switch t {
	case models.PlayerUp:
		return Up
	case models.PlayerDown:
		return Down
	case models.PlayerLeft:
		return Left
	case models.PlayerRight:
		return Right
	}
	return Down
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
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
