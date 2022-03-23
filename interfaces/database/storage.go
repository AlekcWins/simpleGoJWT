package database

import (
	"simpleGoJWT/database/dbModel"
)

type TokenStorage interface {
	Create(token dbModel.Token) (string, error)
	FindOne(id string) (dbModel.Token, error)
	FindOneWithUserID(userUId string) (dbModel.Token, error)
	Update(token dbModel.Token) error
	CreateOrUpdate(token dbModel.Token) (dbModel.Token, error)
	Delete(userUId string) error
}

type UserStorage interface {
	Create(user dbModel.User) (string, error)
	FindAll() (u []dbModel.User, err error)
	FindOne(id string) (dbModel.User, error)
	FindOneWithEmail(id string) (dbModel.User, error)
	Update(user dbModel.User) error
	Delete(id string) error
}
