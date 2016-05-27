package model

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Game struct {
	ID             bson.ObjectId `bson:"_id"`
	FirstPlayer    mgo.DBRef     `bson:"firstPlayer"`
	SecondPlayer   mgo.DBRef     `bson:"secondPlayer,omitempty"`
	StartingPlayer mgo.DBRef     `bson:"startingPlayer,omitempty"`
	Winner         mgo.DBRef     `bson:"winner,omitempty"`
	CreateDate     time.Time     `bson:"createDate"`
	StartDate      time.Time     `bson:"startDate,omitempty"`
	EndDate        time.Time     `bson:"endDate,omitempty"`
	TurnCount      int           `bson:"turnCount"`
}

type Games []Game

func (game Game) DBRef() mgo.DBRef {
	return mgo.DBRef{
		Collection: "game",
		Id:         game.ID,
	}
}

type Turn struct {
	ID              bson.ObjectId `bson:"_id"`
	Game            mgo.DBRef     `bson:"game"`
	Player          mgo.DBRef     `bson:"player"`
	Position        int           `bson:"position"`
	PositionInField int           `bson:"positionInField"`
	Field           int           `bson:"field"`
	NextField       int           `bson:"nextField"`
	RandomField     bool          `bson:"randomField"`
	TurnCount       int           `bson:"turnCount"`
	WonField        bool          `bson:"wonField"`
	WonGame         bool          `bson:"wonGame"`
	CreateDate      time.Time     `bson:"createDate"`
}

type Turns []Turn
