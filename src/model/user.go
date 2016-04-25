package model

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string
}

func (user User) DBRef() mgo.DBRef {
	return mgo.DBRef{
		Collection: "user",
		Id:         user.ID,
	}
}
