package objects

import (
	"bytes"
	"image"
	"log"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const FONT_SIZE = 16
const MESSAGE_WIDTH = 250
const MESSAGE_HEIGHT = 20
const MESSAGE_DURATION = 5 * time.Second

var ChatService = newChat()

type (
	Chat struct {
		mu         sync.Mutex
		font       text.Face
		msgObjects map[string]messageObject
	}
	ChatMessage struct {
		UserID   string `json:"user_id"`
		UserName string `json:"username"`
		Message  string `json:"message"`
	}
	messageObject struct {
		time time.Time
		msg  ChatMessage
		X    int
		Y    int
	}
)

func newChat() *Chat {
	c := &Chat{
		msgObjects: make(map[string]messageObject),
	}
	// load font
	font, err := loadFont(FONT_SIZE)
	if err != nil {
		log.Fatal(err)
	}
	c.font = font

	return c
}

func (c *Chat) AddMessage(msg ChatMessage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// add message to chat
	chatMsg := messageObject{
		time: time.Now(),
		msg:  msg,
	}

	// add message to chat
	c.msgObjects[msg.UserID] = chatMsg
}

func (c *Chat) Update(g IGame, tick uint) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// remove expired messages
	for i, msg := range c.msgObjects {
		if time.Since(msg.time) > MESSAGE_DURATION {
			delete(c.msgObjects, i)
		}
	}

	updateMessages := map[string]messageObject{}

	onlinePlayers := g.OnlinePlayers()
	player := g.Player()
	playerID := player.ID()
	for _, mo := range c.msgObjects {
		var p *Player
		if mo.msg.UserID == playerID {
			p = player
		} else {
			_p, ok := onlinePlayers[mo.msg.UserID]
			if !ok {
				continue
			}
			p = _p
		}

		// truncate message
		lenMsg := len(mo.msg.Message)
		if lenMsg > 24 {
			mo.msg.Message = mo.msg.Message[:24] + "..."
		}

		// get player position
		pos := *p.Position()
		bounds := p.Image().Bounds()
		// center message above player
		center := bounds.Dx() / 2
		mo.X = (pos.X + center) - (len(mo.msg.Message) * 3)
		if len(mo.msg.Message) <= 3 {
			mo.X += len(mo.msg.Message) + 2
		}
		mo.Y = pos.Y - 20

		// update message
		updateMessages[mo.msg.UserID] = mo
	}
	c.msgObjects = updateMessages
	return nil
}

func (c *Chat) Draw(screen *ebiten.Image) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// draw chat objects
	for _, mo := range c.msgObjects {
		// message container
		img := ebiten.NewImageFromImage(
			image.NewRGBA(image.Rect(0, 0, MESSAGE_WIDTH, MESSAGE_HEIGHT)),
		)

		// draw text to container
		// shadow
		shadowOpts := &text.DrawOptions{}
		shadowOpts.ColorScale.Scale(1, 1, 1, 1)
		text.Draw(img, mo.msg.Message, c.font, shadowOpts)
		// text
		textOpts := &text.DrawOptions{}
		textOpts.GeoM.Translate(1, -1)
		textOpts.ColorScale.Scale(0, 0, 0, 1)
		text.Draw(img, mo.msg.Message, c.font, textOpts)

		// draw container to screen
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(mo.X), float64(mo.Y))
		screen.DrawImage(img, opts)
	}
}

func loadFont(size float64) (text.Face, error) {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &text.GoTextFace{
		Source: s,
		Size:   size,
	}, nil
}
