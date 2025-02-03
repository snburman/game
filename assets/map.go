package assets

import "go.mongodb.org/mongo-driver/bson/primitive"

type Portal struct {
	MapID string `json:"map_id" bson:"map_id"`
	X     int    `json:"x" bson:"x"`
	Y     int    `json:"y" bson:"y"`
}

type Map[T any] struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserID   string             `json:"user_id" bson:"user_id"`
	Name     string             `json:"name" bson:"name"`
	Primary  bool               `json:"primary" bson:"primary"`
	Entrance struct {
		X int `json:"x" bson:"x"`
		Y int `json:"y" bson:"y"`
	}
	Portals []Portal `json:"portals" bson:"portals"`
	Data    T        `json:"data" bson:"data"`
}
