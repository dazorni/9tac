package model

import (
	"errors"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	validator "gopkg.in/go-playground/validator.v8"
)

type User struct {
	ID       bson.ObjectId `bson:"_id"`
	Username string        `bson:"username" validate:"required,min=3,max=20,excludesall=!@#$%^&*()_+-=:;?/0x2C"`
}

func (user User) DBRef() mgo.DBRef {
	return mgo.DBRef{
		Collection: "user",
		Id:         user.ID,
	}
}

func (user User) Valid() error {
	validate := validator.New(&validator.Config{TagName: "validate"})

	if err := validate.Struct(user); err != nil {
		return getValidationError(err.(validator.ValidationErrors))
	}

	return nil
}

func getValidationError(allErrors validator.ValidationErrors) error {
	for _, e := range allErrors {
		switch e.Field {
		case "Username":
			if e.Tag == "min" {
				return errors.New("Given username is too short")
			} else if e.Tag == "max" {
				return errors.New("Given username is too long")
			} else if e.Tag == "excludesall" {
				return errors.New("Username contains invalid characters")
			}

			// required
			return errors.New("Username is required")
		default:
			return allErrors
		}
	}

	return errors.New("Something went wrong")
}
