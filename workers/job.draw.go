package workers

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten"
	"github.com/snburman/magicgame/objects"
)

type Draw struct {
	id  string
	cfg DrawConfig
}

type DrawConfig struct {
	Sprite *objects.Object
	Screen *ebiten.Image
	Tick   uint
}

func NewDrawJob(cfg DrawConfig) Job {
	if cfg.Sprite == nil {
		panic("draw job: sprite cannot be nil")
	}

	if cfg.Screen == nil {
		panic("draw job: screen cannot be nil")
	}

	job := NewJob(Draw{
		id:  "draw_job_" + cfg.Sprite.Name() + "_" + uuid.New().String(),
		cfg: cfg,
	})

	return job
}

func (j Draw) ID() string {
	return j.id
}

func (j Draw) Run() error {
	sprite := j.cfg.Sprite
	sprite.Draw(j.cfg.Screen, j.cfg.Tick)
	return nil
}

func (j Draw) Config() DrawConfig {
	return j.cfg
}
