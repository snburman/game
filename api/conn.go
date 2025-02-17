package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/coder/websocket"
	"github.com/snburman/game/config"
)

const (
	ErrReadMessage  = WebsocketError("error reading message")
	ErrParseMessage = WebsocketError("error parsing message")
	ErrWriteMessage = WebsocketError("error writing message")
)

type (
	Conn struct {
		context.Context
		UserID     string
		websocket  *websocket.Conn
		mapService *MapService
		messages   chan []byte
		status     websocket.StatusCode
		close      struct {
			status websocket.StatusCode
			reason string
		}
	}
	WebsocketError string
)

func (e WebsocketError) Error() string {
	return string(e)
}

func NewConn(url string, ms *MapService) (*Conn, error) {
	if ms == nil {
		return nil, errors.New("nil map service")
	}
	if ms.api == nil {
		return nil, errors.New("nil api")
	}
	if ms.api.userID == "" {
		return nil, errors.New("nil user id")
	}
	if url == "" {
		return nil, errors.New("nil url")
	}

	headers := make(map[string][]string)
	headers["CLIENT_ID"] = []string{config.Env().CLIENT_ID}
	headers["CLIENT_SECRET"] = []string{config.Env().CLIENT_SECRET}

	url = url + "/" + ms.api.userID
	ctx := context.Background()
	websocket, res, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		log.Println("res", res)
		return nil, err
	}

	c := &Conn{
		Context:    ctx,
		UserID:     ms.api.userID,
		websocket:  websocket,
		messages:   make(chan []byte, 256),
		mapService: ms,
	}
	go c.listen()
	dispatch := NewDispatch(c, Authenticate, headers)
	dispatch.MarshalAndPublish()

	return c, nil
}

func (c *Conn) listen() {
	go func(c *Conn) {
		defer c.Close()
		var dispatch Dispatch[[]byte]
		for {
			_, b, err := c.websocket.Read(context.Background())
			if err != nil {
				log.Printf("error: %v", err)
				c.close.status = websocket.StatusAbnormalClosure
				c.close.reason = err.Error()
				break
			}
			// Parse dispatch from websocket message
			err = json.Unmarshal(b, &dispatch)
			if err != nil {
				log.Printf("error: %v", err)
				c.close.status = websocket.StatusAbnormalClosure
				c.close.reason = err.Error()
				break
			}

			// Set conn on dispatch
			dispatch.conn = c
			// Route dispatch to appropriate function
			go RouteDispatch(dispatch)
		}
	}(c)

	// outgoing messages
	for {
		msg, ok := <-c.messages
		if !ok || c.websocket == nil {
			break
		}

		if err := c.websocket.Write(c.Context, websocket.MessageText, msg); err != nil {
			log.Println("error writing message", "error", err)
			c.close.status = websocket.StatusAbnormalClosure
			c.close.reason = err.Error()
			break
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
	c.messages <- msg
}

func (c *Conn) Close() error {
	if c == nil {
		return errors.New("cannot close nil connection")
	}
	close(c.messages)
	var status websocket.StatusCode
	if c.status == 0 {
		status = websocket.StatusNormalClosure

	} else {
		status = c.status
	}
	log.Println("websocket connection closing, ", c.UserID)
	return c.websocket.Close(status, c.close.reason)
}
