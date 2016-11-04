package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/dazorni/9tac/src/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type GameStorage struct {
	session      *mgo.Session
	databaseName string
}

func NewGameStorage(session *mgo.Session, databaseName string) *GameStorage {
	return &GameStorage{
		session:      session,
		databaseName: databaseName,
	}
}

func (storage GameStorage) Insert(game *model.Game) error {
	session := storage.session.Copy()
	defer session.Close()

	game.ID = bson.NewObjectId()

	return session.DB(storage.databaseName).C("game").Insert(&game)
}

func (storage GameStorage) JoinGame(game *model.Game, joiningUser model.User) error {
	session := storage.session.Copy()
	defer session.Close()

	err := session.DB(storage.databaseName).C("game").FindId(game.ID).One(&game)

	if err != nil {
		return err
	}

	if game.SecondPlayer.Id != nil {
		return errors.New("There are already two players on this game")
	}

	if game.FirstPlayer.Id == joiningUser.ID {
		return errors.New("A player can not play against himself")
	}

	game.SecondPlayer = joiningUser.DBRef()

	return session.DB(storage.databaseName).C("game").UpdateId(game.ID, &game)
}

func (storage GameStorage) FindAllGamesForUser(user model.User) (model.Games, error) {
	session := storage.session.Copy()
	defer session.Close()

	games := model.Games{}

	query := bson.M{
		"$or": []bson.M{
			bson.M{"firstPlayer.$id": user.ID},
			bson.M{"secondPlayer.$id": user.ID},
		},
	}

	err := session.DB(storage.databaseName).C("game").Find(query).Sort("-startTime").Iter().All(&games)

	return games, err
}

func (storage GameStorage) FindAllOpenGamesForUser(user model.User) (model.Games, error) {
	session := storage.session.Copy()
	defer session.Close()

	games := model.Games{}
	query := bson.M{
		"$or": []bson.M{
			bson.M{"firstPlayer.$id": user.ID},
			bson.M{"secondPlayer.$id": user.ID},
		},
	}

	err := session.DB(storage.databaseName).C("game").Find(query).Sort("-startTime").Iter().All(&games)

	return games, err
}

func (storage GameStorage) FindOne(gameId string) (model.Game, error) {
	session := storage.session.Copy()
	defer session.Close()

	game := model.Game{}

	if bson.IsObjectIdHex(gameId) == false {
		return game, errors.New("Fail")
	}

	objectID := bson.ObjectIdHex(gameId)

	err := session.DB(storage.databaseName).C("game").Find(bson.M{"_id": objectID}).One(&game)

	return game, err
}

func (storage GameStorage) Turn(game *model.Game, player model.User, position int) (model.Turn, error) {
	session := storage.session.Copy()
	defer session.Close()
	defer storage.updateTurnCount(game)

	turn := model.Turn{}
	turn.ID = bson.NewObjectId()
	turn.Game = game.DBRef()
	turn.Player = player.DBRef()
	turn.Position = position
	turn.CreateDate = time.Now()
	turn.NextField = (position % 9 % 3) + ((position / 9 % 3) * 3)
	turn.PositionInField = turn.NextField
	turn.Field = (position % 9 / 3) + ((position / 9 / 3) * 3)

	previousTurn := model.Turn{}

	collection := session.DB(storage.databaseName).C("turn")

	if err := collection.Find(bson.M{"game.$id": game.ID}).Sort("-turnCount").One(&previousTurn); err != nil {
		if err == mgo.ErrNotFound {
			game.StartingPlayer = player.DBRef()
			game.StartDate = time.Now()
			turn.TurnCount = 1
			collection.UpdateId(game.ID, game)

			if err := collection.Insert(turn); err != nil {
				return turn, err
			}

			return turn, nil
		}

		return turn, err
	}

	if previousTurn.Player.Id == turn.Player.Id {
		wrongPlayErr := fmt.Sprintf("Wrong player. Turn: %s", previousTurn.ID)
		return turn, errors.New(wrongPlayErr)
	}

	if previousTurn.RandomField != true && previousTurn.NextField != turn.Field {
		wrongFieldErr := fmt.Sprintf("Wrong field. Previous: %d now: %d, position: %d", previousTurn.NextField, turn.Field, turn.Position)
		return turn, errors.New(wrongFieldErr)
	}

	previousTurns := model.Turns{}

	previousTurnsQuery := bson.M{
		"game.$id": game.ID,
	}

	if err := collection.Find(previousTurnsQuery).Sort("-createDate").Iter().All(&previousTurns); err != nil {
		return turn, err
	}

	for _, previousTurn := range previousTurns {
		if previousTurn.Position == turn.Position {
			sameFieldErr := fmt.Sprintf("Field is already played")
			return turn, errors.New(sameFieldErr)
		}
	}

	if storage.isFieldPlayable(turn.Field, previousTurns) == false {
		fieldAlreadyWonErrMessage := fmt.Sprintf("Field already won")
		return turn, errors.New(fieldAlreadyWonErrMessage)
	}

	var previousFieldPositions [9]bool

	for _, previousTurn := range previousTurns {
		if previousTurn.Field == turn.Field && previousTurn.Player.Id == turn.Player.Id {
			previousFieldPositions[previousTurn.PositionInField] = true
		}
	}

	turn.WonField = checkForWin(previousFieldPositions, turn.PositionInField)

	if turn.WonField == true {
		var previousWonFields [9]bool

		previousWonFieldsQuery := bson.M{
			"game.$id":   game.ID,
			"player.$id": player.ID,
			"wonField":   true,
		}

		winningTurns := model.Turns{}

		if err := collection.Find(previousWonFieldsQuery).Sort("-createDate").Iter().All(&winningTurns); err != nil {
			return turn, err
		}

		for _, turn := range winningTurns {
			previousWonFields[turn.Field] = true
		}

		turn.WonGame = checkForWin(previousWonFields, turn.Field)
	}

	isNextFieldPlayable := true

	if turn.Field == turn.NextField && turn.WonField {
		isNextFieldPlayable = false
	}

	if isNextFieldPlayable == true {
		isNextFieldPlayable = storage.isFieldPlayable(turn.NextField, previousTurns)
	}

	if isNextFieldPlayable == false {
		turn.RandomField = true
	}

	turn.TurnCount = previousTurn.TurnCount + 1
	err := collection.Insert(&turn)

	return turn, err
}

func (storage GameStorage) isFieldPlayable(field int, previousTurns model.Turns) bool {
	isFieldPlayable := true

	for _, previousTurn := range previousTurns {
		if field != previousTurn.Field {
			continue
		}

		if previousTurn.WonField == false {
			continue
		}

		isFieldPlayable = false
	}

	return isFieldPlayable
}

func (storage GameStorage) updateTurnCount(game *model.Game) error {
	session := storage.session.Copy()
	defer session.Close()

	lastTurnQuery := bson.M{"game.$id": game.ID}
	collection := session.DB(storage.databaseName).C("turn")
	lastTurn := model.Turn{}

	if err := collection.Find(lastTurnQuery).Sort("-turnCount").One(&lastTurn); err != nil {
		return err
	}

	game.TurnCount = lastTurn.TurnCount

	return session.DB(storage.databaseName).C("game").UpdateId(game.ID, &game)
}

func checkForWin(previousTurns [9]bool, turnField int) bool {
	combinations := [][]int{
		[]int{0, 4, 8},
		[]int{0, 1, 2},
		[]int{0, 3, 6},
		[]int{1, 4, 7},
		[]int{2, 4, 6},
		[]int{3, 4, 5},
		[]int{2, 5, 8},
		[]int{6, 7, 8},
	}

	for _, combi := range combinations {
		matchingCount := 0

		for _, field := range combi {
			if previousTurns[field] == true || turnField == field {
				matchingCount++
			}
		}

		if matchingCount == 3 {
			return true
		}
	}

	return false
}
