package api

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/snburman/game/config"
)

type (
	Conn struct {
		websocket *websocket.Conn
		ID        string
		LastPing  time.Time
		Messages  chan []byte
	}
)

func NewConn(id string, url string) (*Conn, error) {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	headers := make(map[string][]string)
	headers["CLIENT_ID"] = []string{config.Env().CLIENT_ID}
	headers["CLIENT_SECRET"] = []string{config.Env().CLIENT_SECRET}

	websocket, _, err := dialer.Dial(url, headers)
	if err != nil {
		return nil, err
	}

	c := &Conn{
		websocket: websocket,
		ID:        id,
		Messages:  make(chan []byte, 256),
	}
	return c, nil
}

func (c *Conn) close() error {
	if c == nil {
		return errors.New("cannot close nil connection")
	}
	c.websocket.Close()
	return nil
}

func (c *Conn) listen() {
	go func(c *Conn) {
		defer c.close()
		var dispatch Dispatch[any]
		for {
			_, message, err := c.websocket.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
					websocket.CloseNormalClosure,
				) {
					log.Printf("error: %v", err)
				}
				close(c.Messages)
				break
			}
			// Parse dispatch from websocket message
			err = json.Unmarshal(message, &dispatch)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			// Set conn on dispatch
			dispatch.conn = c
			// Dispatch to message handler
			// handler.in <- dispatch
		}
	}(c)

	// outgoing messages
	for {
		msg, ok := <-c.Messages
		if !ok {
			c.close()
			break
		}
		if c.websocket == nil {
			break
		}

		if err := c.websocket.WriteMessage(1, msg); err != nil {
			log.Println("error writing message", "error", err)
			c.close()
		}
	}
}

func (c *Conn) Publish(msg []byte) {
	// if msg is not json encodable, return
	_, err := json.Marshal(msg)
	if err != nil {
		log.Println("message not json encodable", "error", err)
		return
	}
	if c == nil {
		log.Println("connection severed, message not sent")
		return
	}
	c.Messages <- msg
}

func (c *Conn) Write(p []byte) (n int, err error) {
	c.Messages <- p
	return len(p), nil
}
