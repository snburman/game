package input

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type TouchManager struct {
	mu              sync.Mutex
	allIDs          []ebiten.TouchID
	currentIDs      map[ebiten.TouchID]bool
	newPressedIDs   []ebiten.TouchID
	justPressedIDs  map[ebiten.TouchID]bool
	justReleasedIDs map[ebiten.TouchID]bool
	newReleasedIDs  []ebiten.TouchID
}

func NewTouchManager() *TouchManager {
	return &TouchManager{
		allIDs:          []ebiten.TouchID{},
		currentIDs:      make(map[ebiten.TouchID]bool),
		newPressedIDs:   []ebiten.TouchID{},
		justPressedIDs:  make(map[ebiten.TouchID]bool),
		justReleasedIDs: make(map[ebiten.TouchID]bool),
		newReleasedIDs:  []ebiten.TouchID{},
	}
}

func (tm *TouchManager) CurrentIDs() map[ebiten.TouchID]bool {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.currentIDs
}

func (tm *TouchManager) AllIDs() []ebiten.TouchID {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	return tm.allIDs
}

func (tm *TouchManager) Update() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	tm.justPressedIDs = make(map[ebiten.TouchID]bool)
	tm.justReleasedIDs = make(map[ebiten.TouchID]bool)
	tm.newPressedIDs = []ebiten.TouchID{}
	tm.newReleasedIDs = []ebiten.TouchID{}

	// inpututil.IsTouchJustReleased()

	tm.newPressedIDs = inpututil.AppendJustPressedTouchIDs(tm.newPressedIDs)
	for _, id := range tm.newPressedIDs {
		tm.justPressedIDs[tm.newPressedIDs[id]] = true
		tm.currentIDs[tm.newPressedIDs[id]] = true
		tm.allIDs = append(tm.allIDs, id)
	}

	tm.newReleasedIDs = inpututil.AppendJustReleasedTouchIDs(tm.newReleasedIDs)
	for i := range tm.newReleasedIDs {
		tm.justReleasedIDs[tm.newReleasedIDs[i]] = true
		delete(tm.currentIDs, tm.newReleasedIDs[i])
	}
}
