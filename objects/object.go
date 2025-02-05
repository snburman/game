package objects

import (
	"log"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/input"
)

type ObjectType string

const (
	ObjectTile     ObjectType = "tile"
	ObjectObstacle ObjectType = "obstacle"
	ObjectPlayer   ObjectType = "player"
)

type Object struct {
	name      string
	img       *ebiten.Image
	ObjType   ObjectType
	frames    map[Direction]*ebiten.Image
	position  Position
	direction Direction
	speed     int
	scale     float64
}

type ObjectOptions struct {
	ObjectType ObjectType
	Position   Position
	Direction  Direction
	Speed      int
	Scale      float64
}

type FileImage struct {
	Name string
	Url  string
	Opts ObjectOptions
}

func NewObject(img assets.Image, opts ObjectOptions) *Object {
	scale := 1.0
	if opts.Scale != 0 {
		scale = opts.Scale
	}
	speed := 1
	if opts.Speed != 0 {
		speed = opts.Speed
	}
	var t ObjectType
	switch img.AssetType {
	case assets.Tile:
		t = ObjectTile
	case assets.Object:
		t = ObjectObstacle
	case assets.PlayerUp, assets.PlayerDown, assets.PlayerLeft, assets.PlayerRight:
		t = ObjectPlayer
	}
	if t == "" {
		t = ObjectTile
	}

	o := &Object{
		img:       img.Image,
		name:      img.Name,
		ObjType:   t,
		frames:    map[Direction]*ebiten.Image{},
		position:  opts.Position,
		direction: opts.Direction,
		speed:     speed,
		scale:     scale,
	}

	// assign default name and frames
	if o.ObjType == ObjectPlayer {
		o.name = "player"
		o.frames[Up] = img.Image
		o.frames[Down] = img.Image
		o.frames[Left] = img.Image
		o.frames[Right] = img.Image
	}
	return o
}

func NewObjectFromFile(f FileImage) *Object {
	_img, _, err := ebitenutil.NewImageFromFile("assets/img/" + f.Url)
	if err != nil {
		log.Println("Error loading image", f.Url)
		panic(err)
	}
	img := assets.Image{
		Name:   f.Name,
		Image:  _img,
		Width:  _img.Bounds().Dx(),
		Height: _img.Bounds().Dy(),
		X:      f.Opts.Position.X,
		Y:      f.Opts.Position.Y,
	}

	var t assets.AssetType
	switch f.Opts.ObjectType {
	case ObjectTile:
		t = assets.Tile
	case ObjectObstacle:
		t = assets.Object
	default:
		t = assets.Tile
	}
	img.AssetType = t

	return NewObject(img, f.Opts)
}

func (s *Object) SetFrame(d Direction, img *ebiten.Image) {
	s.frames[d] = img
}

func (s Object) Name() string {
	return s.name
}

func (s Object) Image() *ebiten.Image {
	return s.img
}

func (s Object) Position() Position {
	return s.position
}

func (s *Object) SetPosition(p Position) {
	s.position = p
}

func (s *Object) Direction() Direction {
	return s.direction
}

func (s *Object) SetDirection(d Direction) {
	s.direction = d
}

func (s Object) Speed() int {
	return s.speed
}

func (s *Object) Update(input input.Input, tick uint) error {

	return nil
}

func (s *Object) Draw(screen *ebiten.Image, tick uint) {
	opts := &ebiten.DrawImageOptions{}
	//TODO: check for collision

	opts.GeoM.Scale(float64(s.scale), float64(s.scale))
	opts.GeoM.Translate(float64(s.position.X), float64(s.position.Y))

	// draw player with direction
	if s.ObjType == ObjectPlayer {
		if s.frames[s.direction] != nil {
			screen.DrawImage(s.frames[s.direction], opts)
			return
		} else {
			screen.DrawImage(s.img, opts)
			return
		}
	}
	screen.DrawImage(s.img, opts)
}

type ObjectManager struct {
	mu      sync.Mutex
	Objects []*Objecter
}

func NewObjectManager() *ObjectManager {
	return &ObjectManager{
		Objects: []*Objecter{},
	}
}

func (om *ObjectManager) Add(s Objecter) {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.Objects = append(om.Objects, &s)
}

func (om *ObjectManager) GetAll() []*Objecter {
	om.mu.Lock()
	defer om.mu.Unlock()
	return om.Objects
}

func (om *ObjectManager) RemoveAll() {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.Objects = []*Objecter{}
}
