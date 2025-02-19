package api

import (
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"github.com/snburman/game/models"
	"github.com/snburman/game/objects"
)

const (
	Authenticate        FunctionName = "authenticate"
	LoadOnlinePlayers   FunctionName = "load_online_players"
	LoadNewOnlinePlayer FunctionName = "load_new_online_player"
	RemoveOnlinePlayer  FunctionName = "remove_online_player"
	UpdatePlayer        FunctionName = "update_player"
	Chat                FunctionName = "chat"
)

type (
	FunctionName    string
	Dispatch[T any] struct {
		ID       string       `json:"id"`
		conn     *Conn        `json:"-"`
		Function FunctionName `json:"function"`
		Data     T            `json:"data"`
	}
	PlayerUpdate struct {
		UserID string `json:"user_id"`
		MapID  string `json:"map_id"`
		Dir    int    `json:"dir"`
		Pos    struct {
			X int `json:"x"`
			Y int `json:"y"`
			Z int `json:"z"`
		} `json:"pos"`
	}
)

func NewDispatch[T any](conn *Conn, function FunctionName, data T) Dispatch[T] {
	if conn == nil {
		panic("nil connection")
	}
	if function == "" {
		panic("empty function name")
	}

	return Dispatch[T]{
		ID:       uuid.New().String(),
		conn:     conn,
		Function: function,
		Data:     data,
	}
}

func (d Dispatch[T]) MarshalAndPublish() {
	if d.conn == nil {
		log.Println("nil connection, message not sent")
		return
	}
	databytes, err := json.Marshal(d.Data)
	if err != nil {
		log.Println("dispatch data not json encodable", "error", err)
	}

	dispatch := Dispatch[[]byte]{
		ID:       d.ID,
		conn:     d.conn,
		Function: d.Function,
		Data:     databytes,
	}

	dispatchBytes, err := json.Marshal(dispatch)
	if err != nil {
		log.Println("dispatch struct not json encodable", "error", err)
		return
	}
	d.conn.Publish(dispatchBytes)
}

func ParseDispatch[T any](d Dispatch[[]byte]) Dispatch[T] {
	var dis Dispatch[T]
	err := json.Unmarshal(d.Data, &dis.Data)
	if err != nil {
		log.Println("error parsing dispatch data", "error", err)
		panic(err)
	}
	dis.ID = d.ID
	dis.conn = d.conn
	dis.Function = d.Function
	return dis
}

func RouteDispatch(d Dispatch[[]byte]) {
	if d.conn == nil {
		log.Println("nil connection, message not sent, id", d.ID)
		return
	}

	switch d.Function {
	case LoadOnlinePlayers:
		dispatch := ParseDispatch[[]models.Image](d)
		d.conn.mapService.LoadOnlinePlayers(dispatch.Data)
	case LoadNewOnlinePlayer:
		dispatch := ParseDispatch[[]models.Image](d)
		d.conn.mapService.LoadNewOnlinePlayer(dispatch.Data)
	case UpdatePlayer:
		dispatch := ParseDispatch[PlayerUpdate](d)
		d.conn.mapService.UpdateLocalPlayer(dispatch.Data)
	case RemoveOnlinePlayer:
		dispatch := ParseDispatch[string](d)
		d.conn.mapService.RemoveOnlinePlayerByID(dispatch.Data)
	case Chat:
		dispatch := ParseDispatch[objects.ChatMessage](d)
		objects.ChatService.AddMessage(dispatch.Data)
	}
}
