package input

import "github.com/hajimehoshi/ebiten"

type Keyboard struct {
	keys map[Key]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		keys: make(map[Key]bool),
	}
}

func (k *Keyboard) Press(key Key) {
	k.keys[key] = true
}

func (k *Keyboard) Release(key Key) {
	k.keys[key] = false
}

func (k *Keyboard) IsPressed(key Key) bool {
	return k.keys[key]
}

func (k *Keyboard) Update() {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		k.Press(Up)
	} else {
		k.Release(Up)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		k.Press(Down)
	} else {
		k.Release(Down)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		k.Press(Left)
	} else {
		k.Release(Left)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		k.Press(Right)
	} else {
		k.Release(Right)
	}
}
