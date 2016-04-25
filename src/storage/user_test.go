package storage_test

import (
	"github.com/dazorni/9tac/src/model"
	"labix.org/v2/mgo"

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

	Context("Insert user", func() {
		It("Simple user", func() {
			user := model.User{}
			user.Username = "username"

			err := userStorage.Insert(&user)
			Expect(err).ToNot(HaveOccurred())
			Expect(user.ID.Valid()).To(BeTrue())
		})
	})

	Context("Find by username", func() {
		It("One match", func() {
			username := "username"
			InsertUser(username)

			user, err := userStorage.FindByUsername(username)

			Expect(err).ToNot(HaveOccurred())
			Expect(user.Username).To(Equal(username))
			Expect(user.ID.Valid()).To(BeTrue())
		})

		It("No match", func() {
			username := "nouser"

			_, err := userStorage.FindByUsername(username)

			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(mgo.ErrNotFound))
		})
	})
})
