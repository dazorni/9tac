package storage_test

import (
	"github.com/dazorni/9tac/src/storage"
	"io/ioutil"
	"labix.org/v2/mgo"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var (
	gameStorage     *storage.GameStorage
	userStorage     *storage.UserStorage
	databaseSession *mgo.Session
	databaseName    string
)

var _ = BeforeSuite(func() {
	log.SetOutput(ioutil.Discard)
	databaseName = "ttt_testing"

	databaseSession, err := mgo.Dial("localhost")
	Expect(err).ToNot(HaveOccurred())
	Expect(databaseSession.DB(databaseName).DropDatabase()).ToNot(HaveOccurred())

	gameStorage = storage.NewGameStorage(databaseSession, databaseName)
	userStorage = storage.NewUserStorage(databaseSession, databaseName)
})
