package objects

import (
	"log"
	"sync"

	_ "image/png"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/snburman/game/config"
	"github.com/snburman/game/models"
)

type ObjectType string

const (
	ObjectTile     ObjectType = "tile"
	ObjectObstacle ObjectType = "obstacle"
	ObjectPlayer   ObjectType = "player"
	ObjectPortal   ObjectType = "portal"
)

type Object struct {
	id        string
	name      string
	img       *ebiten.Image
	objType   ObjectType
	frames    map[Direction]*ebiten.Image
	position  Position
	direction Direction
	breached  Breached
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

func NewObject(img models.Image, opts ObjectOptions) *Object {
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
	case models.Tile:
		t = ObjectTile
	case models.Object:
		t = ObjectObstacle
	case models.MapPortal:
		t = ObjectPortal
	case models.PlayerUp, models.PlayerDown, models.PlayerLeft, models.PlayerRight:
		t = ObjectPlayer
	}
	if t == "" {
		t = ObjectTile
	}

	// set id for portal
	// only portals need fixed ids to images
	// other objects have random ids which can be used to specify unique objects
	var id string
	if t == ObjectPortal {
		id = img.ID
	} else {
		id = uuid.New().String()
	}

	o := &Object{
		id:        id,
		img:       img.Image,
		name:      img.Name,
		objType:   t,
		frames:    map[Direction]*ebiten.Image{},
		position:  opts.Position,
		direction: opts.Direction,
		speed:     speed,
		scale:     scale,
	}

	// assign default name and frames
	if o.ObjType() == ObjectPlayer {
		o.name = "player"
		o.frames[Up] = img.Image
		o.frames[Down] = img.Image
		o.frames[Left] = img.Image
		o.frames[Right] = img.Image
	}
	return o
}

func NewObjectFromFile(f FileImage) *Object {
	var _img *ebiten.Image
	var err error
	switch config.Env().ENVIROMENT {
	case "debug":
		_img, _, err = ebitenutil.NewImageFromFile("assets/img/" + f.Url)
	case "build":
		_img, _, err = ebitenutil.NewImageFromFile("../assets/img/" + f.Url)
	}

	if err != nil {
		log.Println("Error loading image", f.Url)
		panic(err)
	}
	img := models.Image{
		Name:   f.Name,
		Image:  _img,
		Width:  _img.Bounds().Dx(),
		Height: _img.Bounds().Dy(),
		X:      f.Opts.Position.X,
		Y:      f.Opts.Position.Y,
	}

	var t models.AssetType
	switch f.Opts.ObjectType {
	case ObjectTile:
		t = models.Tile
	case ObjectObstacle:
		t = models.Object
	default:
		t = models.Tile
	}
	img.AssetType = t

	return NewObject(img, f.Opts)
}

func (o *Object) Update(g IGame, tick uint) error {

	return nil
}

func (o *Object) Draw(screen *ebiten.Image, tick uint) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(float64(o.scale), float64(o.scale))
	opts.GeoM.Translate(float64(o.position.X), float64(o.position.Y))

	// draw player with direction
	if o.ObjType() == ObjectPlayer {
		if o.frames[o.direction] != nil {
			screen.DrawImage(o.frames[o.direction], opts)
			return
		} else {
			screen.DrawImage(o.img, opts)
			return
		}
	}
	screen.DrawImage(o.img, opts)
}

// DetectCollision checks if the object is about to collide with another object
// and sets the breached flags accordingly
func (o *Object) DetectObjectCollision(foreign Objecter) bool {
	// only check for obstacle collision
	if foreign.ObjType() == ObjectTile {
		return false
	}
	// foreign position
	fLeft := float64(foreign.Position().X)
	fRight := float64(foreign.Position().X) + (float64(foreign.Image().Bounds().Dx()) * config.Scale)
	fTop := float64(foreign.Position().Y)
	fBottom := float64(foreign.Position().Y) + (float64(foreign.Image().Bounds().Dy()) * config.Scale)

	// local position
	lLeft := float64(o.position.X)
	lRight := float64(o.position.X) + (float64(o.img.Bounds().Dx()) * config.Scale)
	lTop := float64(o.position.Y)
	lBottom := float64(o.position.Y) + (float64(o.img.Bounds().Dy()) * config.Scale)
	oSpeed := float64(o.speed)

	// approaching from top
	if (lRight >= fLeft && lLeft <= fRight) && (lBottom < fTop) && (lBottom+oSpeed >= fTop) {
		o.breached.Max.Y = true
	}
	// approaching from bottom
	if (lRight >= fLeft && lLeft <= fRight) && (lTop > fBottom) && (lTop-oSpeed <= fBottom) {
		o.breached.Min.Y = true
	}
	// approaching from left
	if (lBottom >= fTop && lTop <= fBottom) && (lRight < fLeft) && (lRight+oSpeed >= fLeft) {
		o.breached.Max.X = true
	}
	// approaching from right
	if (lBottom >= fTop && lTop <= fBottom) && (lLeft > fRight) && (lLeft-oSpeed <= fRight) {
		o.breached.Min.X = true
	}
	return o.breached.Max.X || o.breached.Min.X || o.breached.Max.Y || o.breached.Min.Y
}

func (o *Object) IsCollided(foreign Objecter) bool {
	// only check for obstacle collision
	if foreign.ObjType() == ObjectTile {
		return false
	}
	// foreign position
	fLeft := float64(foreign.Position().X)
	fRight := float64(foreign.Position().X) + (float64(foreign.Image().Bounds().Dx()) * config.Scale)
	fTop := float64(foreign.Position().Y)
	fBottom := float64(foreign.Position().Y) + (float64(foreign.Image().Bounds().Dy()) * config.Scale)

	// local position
	lLeft := float64(o.position.X)
	lRight := float64(o.position.X) + (float64(o.img.Bounds().Dx()) * config.Scale)
	lTop := float64(o.position.Y)
	lBottom := float64(o.position.Y) + (float64(o.img.Bounds().Dy()) * config.Scale)

	// check for collision
	if (lRight >= fLeft && lLeft <= fRight) && (lBottom >= fTop && lTop <= fBottom) {
		return true
	}
	return false
}

// DetectScreenCollision checks if object is about to collide with screen boundaries
func (o *Object) DetectScreenCollision() {
	// approaching from top
	if (float64(o.position.Y) + (float64(o.img.Bounds().Dy()) * config.Scale) + float64(o.speed)) > float64(config.ViewPortHeight) {
		o.breached.Max.Y = true
	} else {
		o.breached.Max.Y = false
	}
	// approaching from bottom
	if (o.position.Y - o.speed) < 0 {
		o.breached.Min.Y = true
	} else {
		o.breached.Min.Y = false
	}
	// approaching from left
	if (float64(o.position.X) + (float64(o.img.Bounds().Dx()) * config.Scale) + float64(o.speed)) > float64(config.ViewPortWidth) {
		o.breached.Max.X = true
	} else {
		o.breached.Max.X = false
	}
	// approaching from right
	if (o.position.X - o.speed) < 0 {
		o.breached.Min.X = true
	} else {
		o.breached.Min.X = false
	}
}

// IsPressed checks if object is pressed by touch at x, y coordinates
func (o *Object) IsPressed(x, y int) bool {
	// touch coordinates
	touchX := float64(x)
	touchY := float64(y)

	// object borders
	top := float64(o.position.Y)
	bottom := float64(o.position.Y) + (float64(o.img.Bounds().Dy()))
	left := float64(o.position.X)
	right := float64(o.position.X) + (float64(o.img.Bounds().Dx()))

	// if touch is within object borders
	if (touchX >= left && touchX <= right) && (touchY >= top && touchY <= bottom) {
		return true
	}
	return false
}

func (o Object) ID() string {
	return o.id
}

func (o Object) Name() string {
	return o.name
}

func (o *Object) SetFrame(d Direction, img *ebiten.Image) {
	o.frames[d] = img
}

func (o Object) ObjType() ObjectType {
	return o.objType
}

func (o Object) Image() *ebiten.Image {
	return o.img
}

func (o *Object) Position() *Position {
	return &o.position
}

func (o *Object) SetPosition(p Position) {
	o.position = p
}

func (o Object) Breached() Breached {
	return o.breached
}

func (o Object) Direction() Direction {
	return o.direction
}

func (o *Object) SetDirection(d Direction) {
	o.direction = d
}

func (o Object) Speed() int {
	return o.speed
}

func (o *Object) SetSpeed(s int) {
	o.speed = s
}

type ObjectManager struct {
	mu      sync.Mutex
	Objects []Objecter
}

func NewObjectManager() *ObjectManager {
	return &ObjectManager{
		Objects: []Objecter{},
	}
}

func (om *ObjectManager) Add(s Objecter) {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.Objects = append(om.Objects, s)
}

func (om *ObjectManager) Get(name string) Objecter {
	om.mu.Lock()
	defer om.mu.Unlock()
	for _, o := range om.Objects {
		if (o).Name() == name {
			return o
		}
	}
	return nil
}

func (om *ObjectManager) GetAll() []Objecter {
	om.mu.Lock()
	defer om.mu.Unlock()
	return om.Objects
}

func (om *ObjectManager) SetAll(objs []Objecter) {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.Objects = objs
}

func (om *ObjectManager) RemoveAll() {
	om.mu.Lock()
	defer om.mu.Unlock()
	om.Objects = []Objecter{}
}
