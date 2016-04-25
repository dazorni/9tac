package storage_test

import (
	"github.com/dazorni/9tac/src/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {
	InsertUser := func(username string) model.User {
		user := model.User{}
		user.Username = username

		err := userStorage.Insert(&user)
		Expect(err).ToNot(HaveOccurred())

		return user
	}

	InsertGame := func(player model.User) model.Game {
		game := model.Game{}
		game.FirstPlayer = player.DBRef()

		err := gameStorage.Insert(&game)
		Expect(err).ToNot(HaveOccurred())

		return game
	}

	InsertFullGame := func(playerOne model.User, playerTwo model.User) model.Game {
		game := model.Game{}
		game.FirstPlayer = playerOne.DBRef()

		Expect(gameStorage.Insert(&game)).ToNot(HaveOccurred())
		Expect(gameStorage.JoinGame(&game, playerTwo)).ToNot(HaveOccurred())

		return game
	}

	Context("Insert game", func() {
		It("Simple game", func() {
			user := InsertUser("startingplayer")

			game := model.Game{}
			game.FirstPlayer = user.DBRef()

			err := gameStorage.Insert(&game)
			Expect(err).ToNot(HaveOccurred())
			Expect(game.ID.Valid()).To(BeTrue())
		})
	})

	Context("Join user to game", func() {
		It("User join to game", func() {
			user := InsertUser("first")
			joiner := InsertUser("joiner")
			game := InsertGame(user)

			err := gameStorage.JoinGame(&game, joiner)
			Expect(err).ToNot(HaveOccurred())
			Expect(game.SecondPlayer).ToNot(BeNil())
		})

		It("No game found", func() {
			joiner := InsertUser("joiner")
			game := model.Game{}
			game.ID = bson.NewObjectId()

			err := gameStorage.JoinGame(&game, joiner)
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mgo.ErrNotFound))
		})

		It("Same user does not work", func() {
			user := InsertUser("same")
			game := InsertGame(user)

			err := gameStorage.JoinGame(&game, user)
			Expect(err).To(HaveOccurred())
			Expect(game.FirstPlayer.Id).ToNot(Equal(game.SecondPlayer.Id))
			Expect(err.Error()).To(ContainSubstring("A player can not play against himself"))
		})
	})

	Context("Games for user", func() {
		It("Match two games", func() {
			user := InsertUser("username")
			InsertGame(user)
			InsertGame(user)

			games, err := gameStorage.FindAllGamesForUser(user)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(games)).To(Equal(2))
		})

		It("No match", func() {
			user := InsertUser("username")
			games, err := gameStorage.FindAllGamesForUser(user)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(games)).To(Equal(0))
		})
	})

	Context("Open games for user", func() {
		It("Simple find", func() {
			user := InsertUser("progamer")
			secondUser := InsertUser("noob")

			InsertGame(user)
			InsertGame(user)
			InsertGame(secondUser)

			openGames, err := gameStorage.FindAllOpenGamesForUser(user)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(openGames)).To(Equal(2))
		})

		It("No match", func() {
			user := InsertUser("test")
			openGames, err := gameStorage.FindAllOpenGamesForUser(user)

			Expect(err).ToNot(HaveOccurred())
			Expect(len(openGames)).To(Equal(0))
		})
	})

	Context("Turn", func() {
		It("First turn", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)
			startDate := time.Now()

			turn, err := gameStorage.Turn(&game, playerOne, 27)
			Expect(err).ToNot(HaveOccurred())
			Expect(game.StartDate).To(BeTemporally(">", startDate))
			Expect(game.StartingPlayer.Id).To(Equal(playerOne.DBRef().Id))
			Expect(turn.NextField).To(Equal(0))
		})

		It("Insert one complete game", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			var testTurns = []struct {
				Player    model.User
				X         int
				Y         int
				Field     int
				NextField int
				WonField  bool
				WonGame   bool
				lastTurn  bool
			}{
				{playerTwo, 6, 2, 2, 6, false, false, false},
				{playerOne, 2, 7, 6, 5, false, false, false},
				{playerTwo, 6, 3, 5, 0, false, false, false},
				{playerOne, 2, 2, 0, 8, false, false, false},
				{playerTwo, 6, 6, 8, 0, false, false, false},
				{playerOne, 1, 2, 0, 7, false, false, false},
				{playerTwo, 5, 8, 7, 8, false, false, false},
				{playerOne, 6, 7, 8, 3, false, false, false},
				{playerTwo, 1, 5, 3, 7, false, false, false},
				{playerOne, 4, 6, 7, 1, false, false, false},
				{playerTwo, 5, 0, 1, 2, false, false, false},
				{playerOne, 8, 0, 2, 2, false, false, false},
				{playerTwo, 7, 0, 2, 1, false, false, false},
				{playerOne, 4, 0, 1, 1, false, false, false},
				{playerTwo, 3, 1, 1, 3, false, false, false},
				{playerOne, 0, 5, 3, 6, false, false, false},
				{playerTwo, 0, 8, 6, 6, false, false, false},
				{playerOne, 2, 8, 6, 8, false, false, false},
				{playerTwo, 8, 8, 8, 8, false, false, false},
				{playerOne, 6, 8, 8, 6, false, false, false},
				{playerTwo, 1, 6, 6, 1, false, false, false},
				{playerOne, 4, 2, 1, 7, false, false, false},
				{playerTwo, 4, 8, 7, 7, false, false, false},
				{playerOne, 3, 8, 7, 6, false, false, false},
				{playerTwo, 0, 6, 6, 0, false, false, false},
				{playerOne, 2, 1, 0, 5, false, false, false},
				{playerTwo, 8, 5, 5, 8, false, false, false},
				{playerOne, 7, 7, 8, 4, false, false, false},
				{playerTwo, 4, 4, 4, 4, false, false, false},
				{playerOne, 3, 4, 4, 3, false, false, false},
				{playerTwo, 1, 4, 3, 4, false, false, false},
				{playerOne, 3, 3, 4, 0, false, false, false},
				{playerTwo, 0, 1, 0, 3, false, false, false},
				{playerOne, 0, 3, 3, 0, false, false, false},
				{playerTwo, 0, 0, 0, 0, false, false, false},
				{playerOne, 2, 0, 0, 2, true, false, false},
				{playerTwo, 8, 1, 2, 5, false, false, false},
				{playerOne, 7, 5, 5, 7, false, false, false},
				{playerTwo, 3, 6, 7, 0, false, false, false},
				{playerOne, 0, 2, 0, 6, false, false, false},
				{playerTwo, 2, 6, 6, 2, true, false, false},
				{playerOne, 6, 0, 2, 0, false, false, false},
				{playerTwo, 1, 1, 0, 4, false, false, false},
				{playerOne, 3, 5, 4, 6, true, false, false},
				{playerTwo, 1, 7, 6, 4, false, false, false},
				{playerOne, 4, 3, 4, 1, false, false, false},
				{playerTwo, 3, 0, 1, 0, false, false, false},
				{playerOne, 1, 0, 0, 1, false, false, false},
				{playerTwo, 4, 1, 1, 4, false, false, false},
				{playerOne, 5, 5, 4, 8, false, false, false},
				{playerTwo, 7, 6, 8, 1, false, false, false},
				{playerOne, 5, 2, 1, 8, false, false, false},
				{playerTwo, 8, 6, 8, 2, true, false, false},
				{playerOne, 8, 2, 2, 8, false, false, false},
				{playerTwo, 8, 7, 8, 5, false, false, false},
				{playerOne, 7, 4, 5, 4, false, false, false},
				{playerTwo, 5, 4, 4, 5, false, false, false},
				{playerOne, 6, 5, 5, 6, false, false, false},
				{playerTwo, 0, 7, 6, 3, false, false, false},
				{playerOne, 2, 3, 3, 2, false, false, false},
				{playerTwo, 7, 1, 2, 4, false, false, false},
				{playerOne, 5, 3, 4, 2, false, false, false},
				{playerTwo, 7, 2, 2, 7, true, false, false},
				{playerOne, 4, 7, 7, 4, false, false, false},
				{playerTwo, 4, 5, 4, 7, false, false, false},
				{playerOne, 5, 6, 7, 2, true, false, false},
				{playerTwo, 6, 1, 2, 3, false, false, false},
				{playerOne, 0, 4, 3, 3, true, false, false},
				{playerTwo, 2, 5, 3, 8, false, false, false},
				{playerOne, 7, 8, 8, 7, false, false, false},
				{playerTwo, 3, 7, 7, 3, false, false, false},
				{playerOne, 1, 3, 3, 1, false, false, false},
				{playerTwo, 3, 2, 1, 6, true, false, false},
				{playerOne, 1, 8, 6, 7, false, false, false},
				{playerTwo, 5, 7, 7, 5, false, false, false},
				{playerOne, 8, 3, 5, 2, true, true, true},
			}

			for _, turn := range testTurns {
				position := turn.X + (turn.Y * 9)

				resultTurn, err := gameStorage.Turn(&game, turn.Player, position)

				Expect(err).ToNot(HaveOccurred())
				Expect(resultTurn.Field).To(Equal(turn.Field))
				Expect(resultTurn.WonField).To(Equal(turn.WonField))
				Expect(resultTurn.WonGame).To(Equal(turn.WonGame))
				Expect(resultTurn.Position).To(Equal(position))
				Expect(resultTurn.NextField).To(Equal(turn.NextField))
			}
		})

		It("Insert invalid turn", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			gameStorage.Turn(&game, playerOne, 7)
			_, err := gameStorage.Turn(&game, playerOne, 7)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Wrong player"))
		})

		It("Invalid Field", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			gameStorage.Turn(&game, playerOne, 7)
			_, err := gameStorage.Turn(&game, playerTwo, 24)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Wrong field"))
		})

		It("Turn with already played field", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			gameStorage.Turn(&game, playerOne, 80)
			_, err := gameStorage.Turn(&game, playerTwo, 80)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Field is already played"))
		})
	})
})
