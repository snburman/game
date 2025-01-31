package objects

import (
	"sync"

	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/game/assets"
	"github.com/snburman/game/input"
)

type ObjectName string

type Object struct {
	name         string
	img          *ebiten.Image
	currentFrame int
	position     Position
	direction    Direction
	speed        int
	scale        float64
}

type ObjectOptions struct {
	Position  Position
	Direction Direction
	Speed     int
	Scale     float64
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
	return &Object{
		img:  img.Image,
		name: img.Name,
		// frames:    img.Frames,
		position:  opts.Position,
		direction: opts.Direction,
		speed:     speed,
		scale:     scale,
	}
}

func (s Object) Name() string {
	return s.name
}

func (s Object) Image() *ebiten.Image {
	return s.img
}

func (s Object) CurrentFrame() int {
	return s.currentFrame
}

func (s *Object) SetCurrentFrame(f int) {
	s.currentFrame = f
}

func (s Object) Position() Position {
	return s.position
}

func (s *Object) SetPosition(p Position) {
	s.position = p
}

func (s Object) Direction() Direction {
	return s.direction
}

func (s *Object) SetDirection(d Direction) {
	s.direction = d
}

func (s Object) Speed() int {
	return s.speed
}

func (s *Object) Update(screen *ebiten.Image, input input.Input, tick uint) error {

	return nil
}

func (s *Object) Draw(screen *ebiten.Image, tick uint) {
	opts := &ebiten.DrawImageOptions{}
	//TODO: check for collision

	opts.GeoM.Scale(float64(s.scale), float64(s.scale))
	opts.GeoM.Translate(float64(s.position.X), float64(s.position.Y))

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
