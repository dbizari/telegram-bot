package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type GameSession struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	OwnerId string             `json:"owner_id" bson:"owner_id"`
	Users   []UserInfo         `json:"users" bson:"users"`
	Status  string             `json:"status" bson:"status"`
}

type UserInfo struct {
	UserId string `json:"user_id" bson:"user_id"`
	Role   string `json:"role" bson:"role"`
	Alive  bool   `json:"alive" bson:"alive"`
}
