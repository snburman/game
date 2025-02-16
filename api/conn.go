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
		UserID     string
		websocket  *websocket.Conn
		mapService *MapService
		messages   chan []byte
	}
)

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

	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}
	headers := make(map[string][]string)
	headers["CLIENT_ID"] = []string{config.Env().CLIENT_ID}
	headers["CLIENT_SECRET"] = []string{config.Env().CLIENT_SECRET}

	url = url + "/" + ms.api.userID
	log.Println("connecting to websocket", "url", url)
	websocket, res, err := dialer.Dial(url, headers)
	if err != nil {
		log.Println("res", res)
		return nil, err
	}
	log.Println("websocket connection established", "status", res.Status)

	c := &Conn{
		UserID:     ms.api.userID,
		websocket:  websocket,
		messages:   make(chan []byte, 256),
		mapService: ms,
	}
	go c.listen()
	return c, nil
}

func (c *Conn) listen() {
	go func(c *Conn) {
		defer c.Close()
		var dispatch Dispatch[[]byte]
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
				close(c.messages)
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
			// Route dispatch to appropriate function
			RouteDispatch(dispatch)
		}
	}(c)

	// outgoing messages
	for {
		msg, ok := <-c.messages
		if !ok {
			c.Close()
			break
		}
		if c.websocket == nil {
			break
		}

		if err := c.websocket.WriteMessage(1, msg); err != nil {
			log.Println("error writing message", "error", err)
			c.Close()
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
	c.websocket.Close()
	log.Println("websocket connection closed, ", c.UserID)
	return nil
}
