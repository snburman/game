package input

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Touch struct {
	X, Y int
}

type TouchManager struct {
	touchIDs []ebiten.TouchID
	touches  map[ebiten.TouchID]Touch
}

func NewTouchManager() *TouchManager {
	return &TouchManager{
		touchIDs: []ebiten.TouchID{},
		touches:  make(map[ebiten.TouchID]Touch),
	}
}

func (tm *TouchManager) Touches() map[ebiten.TouchID]Touch {
	return tm.touches
}

func (tm *TouchManager) TouchIDs() []ebiten.TouchID {
	return tm.touchIDs
}

func (tm *TouchManager) Update() {
	for id := range tm.touches {
		if inpututil.IsTouchJustReleased(id) {
			fmt.Println("Touches", tm.touches)
			fmt.Println("Touch released:", id)
			delete(tm.touches, id)
		}
	}

	tm.touchIDs = inpututil.AppendJustPressedTouchIDs(tm.touchIDs)
	for _, id := range tm.touchIDs {
		x, y := ebiten.TouchPosition(id)
		tm.touches[id] = Touch{X: x, Y: y}
	}
}
