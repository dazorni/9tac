package storage

import (
	"github.com/dazorni/9tac/src/model"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type UserStorage struct {
	session      *mgo.Session
	databaseName string
}

func NewUserStorage(session *mgo.Session, databaseName string) *UserStorage {
	return &UserStorage{
		session:      session,
		databaseName: databaseName,
	}
}

func (storage UserStorage) Insert(user *model.User) error {
	session := storage.session.Copy()
	defer session.Close()

	user.ID = bson.NewObjectId()

	// TODO: Do validation

	return session.DB(storage.databaseName).C("user").Insert(&user)
}

func (storage UserStorage) FindByUsername(username string) (model.User, error) {
	session := storage.session.Copy()
	defer session.Close()

	user := model.User{}
	err := session.DB(storage.databaseName).C("user").Find(bson.M{"username": username}).One(&user)

	return user, err
}

func (storage UserStorage) FindByRef(ref mgo.DBRef) (model.User, error) {
	session := storage.session.Copy()
	defer session.Close()

	user := model.User{}
	err := session.DB(storage.databaseName).FindRef(&ref).One(&user)

	return user, err
}
