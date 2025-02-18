package input

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Touch struct {
	X, Y int
}

type TouchManager struct {
	touchIDs map[ebiten.TouchID]bool
}

func NewTouchManager() *TouchManager {
	return &TouchManager{
		touchIDs: make(map[ebiten.TouchID]bool, 128),
	}
}

func (tm *TouchManager) TouchIDs() map[ebiten.TouchID]bool {
	return tm.touchIDs
}

func (tm *TouchManager) Update() {
	var justReleasedIDs []ebiten.TouchID
	justReleasedIDs = inpututil.AppendJustReleasedTouchIDs(justReleasedIDs)
	for _, id := range justReleasedIDs {
		delete(tm.touchIDs, id)
	}
	newIDs := make([]ebiten.TouchID, 0, 56)
	newIDs = inpututil.AppendJustPressedTouchIDs(newIDs)
	for _, id := range newIDs {
		tm.touchIDs[id] = true
	}
}

func (tm *TouchManager) Touches() []Touch {
	touches := make([]Touch, 0, len(tm.touchIDs))
	for id := range tm.touchIDs {
		x, y := ebiten.TouchPosition(id)
		touches = append(touches, Touch{X: x, Y: y})
	}
	return touches
}
