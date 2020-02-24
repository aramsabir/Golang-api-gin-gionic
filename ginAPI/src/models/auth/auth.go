package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

//User
type Auth struct {
	ID        bson.ObjectId `json:"id,omitempty" bson:"_id,omitempty"`
	Token     string        `bson:"token,omitempty"`
	UserId    bson.ObjectId `bson:"user_id,omitempty"`
	CreatedAt time.Time     `bson:"created_at"`
}

type Auths []Auth
