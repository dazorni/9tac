package model_test

import (
	"github.com/dazorni/9tac/src/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Model", func() {
	Context("Validate user", func() {
		It("Valid username", func() {
			user := model.User{}
			user.Username = "username"

			err := user.Valid()
			Expect(err).ToNot(HaveOccurred())
		})

		It("Too long username", func() {
			user := model.User{}
			user.Username = "morethantwentycharactersusername"

			err := user.Valid()
			Expect(err).To(HaveOccurred())
		})

		It("Too short username", func() {
			user := model.User{}
			user.Username = "la"

			err := user.Valid()
			Expect(err).To(HaveOccurred())
		})

		It("Invalid characters", func() {
			user := model.User{}
			user.Username = "invalid!"

			err := user.Valid()
			Expect(err).To(HaveOccurred())
		})
	})
})
