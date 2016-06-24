package storage_test

import (
	"time"

	"github.com/dazorni/9tac/src/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

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

		It("Second player is already set", func() {
			firstPlayer := InsertUser("first")
			secondPlayer := InsertUser("second")
			thirdPlayer := InsertUser("third")

			game := InsertGame(firstPlayer)
			gameStorage.JoinGame(&game, secondPlayer)

			err := gameStorage.JoinGame(&game, thirdPlayer)

			Expect(err).To(HaveOccurred())
			Expect(game.SecondPlayer.Id).ToNot(Equal(thirdPlayer.ID))
			Expect(err.Error()).To(ContainSubstring("There are already two players on this game"))
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

		It("Check turn count", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			turn, err := gameStorage.Turn(&game, playerOne, 27)
			Expect(err).ToNot(HaveOccurred())
			Expect(turn.TurnCount).To(Equal(1))
			Expect(game.TurnCount).To(Equal(1))

			databaseGame, err := gameStorage.FindOne(game.ID.Hex())
			Expect(err).ToNot(HaveOccurred())
			Expect(databaseGame.TurnCount).To(Equal(1))
		})

		It("Insert one complete game", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			var testTurns = []struct {
				Player            model.User
				Position          int
				Field             int
				NextField         int
				WonField          bool
				WonGame           bool
				IsNextFieldRandom bool
				lastTurn          bool
			}{
				{playerTwo, 77, 7, 8, false, false, false, false},
				{playerOne, 62, 8, 2, false, false, false, false},
				{playerTwo, 8, 2, 2, false, false, false, false},
				{playerOne, 24, 2, 6, false, false, false, false},
				{playerTwo, 72, 6, 6, false, false, false, false},
				{playerOne, 54, 6, 0, false, false, false, false},
				{playerTwo, 0, 0, 0, false, false, false, false},
				{playerOne, 1, 0, 1, false, false, false, false},
				{playerTwo, 3, 1, 0, false, false, false, false},
				{playerOne, 19, 0, 7, false, false, false, false},
				{playerTwo, 76, 7, 7, false, false, false, false},
				{playerOne, 66, 7, 3, false, false, false, false},
				{playerTwo, 38, 3, 5, false, false, false, false},
				{playerOne, 44, 5, 5, false, false, false, false},
				{playerTwo, 35, 5, 2, false, false, false, false},
				{playerOne, 16, 2, 4, false, false, false, false},
				{playerTwo, 40, 4, 4, false, false, false, false},
				{playerOne, 41, 4, 5, false, false, false, false},
				{playerTwo, 43, 5, 4, false, false, false, false},
				{playerOne, 50, 4, 8, false, false, false, false},
				{playerTwo, 61, 8, 1, false, false, false, false},
				{playerOne, 4, 1, 1, false, false, false, false},
				{playerTwo, 5, 1, 2, false, false, false, false},
				{playerOne, 25, 2, 7, false, false, false, false},
				{playerTwo, 75, 7, 6, true, false, false, false},
				{playerOne, 74, 6, 8, false, false, false, false},
				{playerTwo, 80, 8, 8, false, false, false, false},
				{playerOne, 78, 8, 6, false, false, false, false},
				{playerTwo, 65, 6, 5, false, false, false, false},
				{playerOne, 51, 5, 6, false, false, false, false},
				{playerTwo, 63, 6, 3, false, false, false, false},
				{playerOne, 45, 3, 6, false, false, false, false},
				{playerTwo, 64, 6, 4, true, false, false, false},
				{playerOne, 32, 4, 2, true, false, false, false},
				{playerTwo, 7, 2, 1, false, false, false, false},
				{playerOne, 14, 1, 5, false, false, false, false},
				{playerTwo, 33, 5, 0, false, false, false, false},
				{playerOne, 10, 0, 4, true, false, true, false},
				{playerTwo, 70, 8, 4, false, false, true, false},
				{playerOne, 69, 8, 3, false, false, false, false},
				{playerTwo, 28, 3, 1, false, false, false, false},
				{playerOne, 12, 1, 3, false, false, false, false},
				{playerTwo, 29, 3, 2, false, false, false, false},
				{playerOne, 15, 2, 3, false, false, false, false},
				{playerTwo, 36, 3, 3, false, false, false, false},
				{playerOne, 27, 3, 0, false, false, true, false},
				{playerTwo, 60, 8, 0, true, true, true, true},
			}

			for _, turn := range testTurns {
				resultTurn, err := gameStorage.Turn(&game, turn.Player, turn.Position)

				Expect(err).ToNot(HaveOccurred())
				Expect(resultTurn.Field).To(Equal(turn.Field))
				Expect(resultTurn.WonField).To(Equal(turn.WonField))
				Expect(resultTurn.WonGame).To(Equal(turn.WonGame))
				Expect(resultTurn.Position).To(Equal(turn.Position))
				Expect(resultTurn.RandomField).To(Equal(turn.IsNextFieldRandom))
				Expect(resultTurn.NextField).To(Equal(turn.NextField))
			}
		})

		It("Insert invalid turn in won field", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			gameStorage.Turn(&game, playerTwo, 75)
			gameStorage.Turn(&game, playerOne, 73)
			gameStorage.Turn(&game, playerTwo, 66)
			gameStorage.Turn(&game, playerOne, 46)
			gameStorage.Turn(&game, playerTwo, 57)
			_, err := gameStorage.Turn(&game, playerTwo, 19)

			Expect(err).To(HaveOccurred())
		})

		It("Insert game with random field", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			var testTurns = []struct {
				Player            model.User
				Position          int
				Field             int
				NextField         int
				WonField          bool
				WonGame           bool
				IsNextFieldRandom bool
			}{
				{playerTwo, 75, 7, 6, false, false, false},
				{playerOne, 73, 6, 7, false, false, false},
				{playerTwo, 77, 7, 8, false, false, false},
				{playerOne, 79, 8, 7, false, false, false},
				{playerTwo, 76, 7, 7, true, false, true},
			}

			for _, turn := range testTurns {
				resultTurn, err := gameStorage.Turn(&game, turn.Player, turn.Position)

				Expect(err).ToNot(HaveOccurred())
				Expect(resultTurn.Field).To(Equal(turn.Field))
				Expect(resultTurn.WonField).To(Equal(turn.WonField))
				Expect(resultTurn.WonGame).To(Equal(turn.WonGame))
				Expect(resultTurn.Position).To(Equal(turn.Position))
				Expect(resultTurn.NextField).To(Equal(turn.NextField))
				Expect(resultTurn.RandomField).To(Equal(turn.IsNextFieldRandom))
			}
		})

		It("Insert game with random field at the end", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			var testTurns = []struct {
				Player            model.User
				Position          int
				Field             int
				NextField         int
				WonField          bool
				WonGame           bool
				IsNextFieldRandom bool
			}{
				{playerTwo, 75, 7, 6, false, false, false},
				{playerOne, 73, 6, 7, false, false, false},
				{playerTwo, 77, 7, 8, false, false, false},
				{playerOne, 79, 8, 7, false, false, false},
				{playerTwo, 76, 7, 7, true, false, true},
			}

			for _, turn := range testTurns {
				resultTurn, err := gameStorage.Turn(&game, turn.Player, turn.Position)

				Expect(err).ToNot(HaveOccurred())
				Expect(resultTurn.Field).To(Equal(turn.Field))
				Expect(resultTurn.WonField).To(Equal(turn.WonField))
				Expect(resultTurn.WonGame).To(Equal(turn.WonGame))
				Expect(resultTurn.Position).To(Equal(turn.Position))
				Expect(resultTurn.NextField).To(Equal(turn.NextField))
				Expect(resultTurn.RandomField).To(Equal(turn.IsNextFieldRandom))
			}
		})

		It("Insert invalid turn in already won field", func() {
			playerOne := InsertUser("playerOne")
			playerTwo := InsertUser("playerTwo")
			game := InsertFullGame(playerOne, playerTwo)

			gameStorage.Turn(&game, playerOne, 3)
			gameStorage.Turn(&game, playerTwo, 1)
			gameStorage.Turn(&game, playerOne, 5)
			gameStorage.Turn(&game, playerTwo, 7)
			gameStorage.Turn(&game, playerOne, 4)
			_, err := gameStorage.Turn(&game, playerTwo, 12)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Field already won"))
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
