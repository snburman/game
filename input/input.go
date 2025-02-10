package input

type Key int

const (
	Up Key = iota
	Down
	Left
	Right
	Space
)

type Input interface {
	Press(Key)
	Release(Key)
	IsPressed(Key) bool
}

type InputFunctions map[Key]func()
