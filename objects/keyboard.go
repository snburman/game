package objects

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/snburman/game/config"
	"github.com/snburman/game/input"
)

type Keyboard struct {
	keys map[input.Key]bool
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		keys: make(map[input.Key]bool),
	}
}

func (k *Keyboard) Press(key input.Key) {
	k.keys[key] = true
}

func (k *Keyboard) Release(key input.Key) {
	k.keys[key] = false
}

func (k *Keyboard) IsPressed(key input.Key) bool {
	return k.keys[key]
}

func (k *Keyboard) Update(g IGame) {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player().SetSpeed(config.RunSpeed)
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		k.Press(input.Up)
	} else {
		k.Release(input.Up)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		k.Press(input.Down)
	} else {
		k.Release(input.Down)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		k.Press(input.Left)
	} else {
		k.Release(input.Left)
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		k.Press(input.Right)
	} else {
		k.Release(input.Right)
	}
}
