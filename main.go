package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/dazorni/9tac/src/model"
	"github.com/dazorni/9tac/src/storage"

	"github.com/googollee/go-socket.io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func main() {
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	if "" == os.Getenv("PORT") {
		os.Setenv("PORT", "5000")
	}

	if "" == os.Getenv("DB_NAME") {
		os.Setenv("DB_NAME", "9tac")
	}

	if "" == os.Getenv("DB_HOST") {
		os.Setenv("DB_HOST", "localhost:27017")
	}

	if "" == os.Getenv("PUBLIC_DIR") {
		os.Setenv("PUBLIC_DIR", "public")
	}

	connectionString := "localhost:27017"

	if os.Getenv("MONGO_PORT_27017_TCP_ADDR") != "" && os.Getenv("MONGO_PORT_27017_TCP_PORT") != "" {
		connectionString = fmt.Sprintf("%s:%s", os.Getenv("MONGO_PORT_27017_TCP_ADDR"), os.Getenv("MONGO_PORT_27017_TCP_PORT"))
	}

	session, databaseError := mgo.Dial(connectionString)

	if databaseError != nil {
		log.Fatal(databaseError)
	}

	database := os.Getenv("DB_NAME")

	gameStorage := storage.NewGameStorage(session, database)
	userStorage := storage.NewUserStorage(session, database)

	server.On("connection", func(socket socketio.Socket) {
		log.Println("user connected")
		socket.Join("game")

		socket.On("user", func(username string) {
			user := model.User{}
			user.Username = username

			if err := userStorage.Insert(&user); err != nil {
				log.Print(err)

				return
			}

			socket.Emit("user", user)
		})

		socket.On("game:new", func(username string) {
			// TODO: Check if player is correct

			log.Printf("game:new(username: %s)", username)

			user, err := userStorage.FindByUsername(username)

			if err != nil {
				log.Printf("game:new(username: %s): couldn't find user", username)

				user.Username = username

				if err := userStorage.Insert(&user); err != nil {
					log.Printf("game:new(username: %s): Couldn't create new user", username)

					return
				}

				log.Printf("game:new(username: %s): Created new user", username)
			}

			game := model.Game{}
			game.FirstPlayer = user.DBRef()

			if err := gameStorage.Insert(&game); err != nil {
				log.Print(err)

				return
			}

			socket.Join(game.ID.String())
			socket.Emit("game:new", game.ID)

			log.Printf("game:new(username: %s): Game started (%s)", username, game.ID)
		})

		socket.On("game:join", func(username string, gameCode string) {
			// do not user id for game code
			if bson.IsObjectIdHex(gameCode) != true {
				log.Printf("'%s' is no hex code", gameCode)
				return
			}

			user, err := userStorage.FindByUsername(username)

			if err != nil {
				user.Username = username

				if err := userStorage.Insert(&user); err != nil {
					log.Print(err)

					return
				}
			}

			game := model.Game{}
			game.ID = bson.ObjectIdHex(gameCode)

			if err := gameStorage.JoinGame(&game, user); err != nil {
				// emit error to the user
				log.Print(err)

				return
			}

			secondPlayer, findErr := userStorage.FindByRef(game.FirstPlayer)

			if findErr != nil {
				log.Printf("game:join(%s, %s): Couldn't find second user", username, gameCode)
			}

			startingField := rand.Intn(8)

			socket.Join(game.ID.String())
			socket.Emit("game:join", username, secondPlayer.Username, gameCode, startingField)
			socket.BroadcastTo(game.ID.String(), "game:join", username, secondPlayer.Username, gameCode, startingField)

			log.Printf("game:join by '%s' to game '%s'", username, game.ID.String())
		})

		socket.On("game:turn", func(username string, gameCode string, position int) {
			log.Printf("game:turn(%s, %s, %d)", username, gameCode, position)

			user, err := userStorage.FindByUsername(username)

			if err != nil {
				log.Printf("game:turn(%s, %s, %d): Couldn't find user", username, gameCode, position)

				return
			}

			game, findErr := gameStorage.FindOne(gameCode)

			if findErr != nil {
				log.Printf("game:turn(%s, %s, %d): Couldn't find game", username, gameCode, position)

				return
			}

			turn, turnErr := gameStorage.Turn(&game, user, position)

			if turnErr != nil {
				log.Printf("game:turn(%s, %s, %d): Couldn't insert turn", username, gameCode, position)

				return
			}

			socket.Emit("game:turn:draw", turn, username)
			socket.BroadcastTo(game.ID.String(), "game:turn:draw", turn, username)
		})

		socket.On("disconnection", func() {
			// TODO: change state for user on pending
			log.Println("user disconnect")
		})
	})

	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)

	http.Handle("/", http.FileServer(http.Dir(os.Getenv("PUBLIC_DIR"))))

	log.Printf("Serving at localhost:%s...", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
