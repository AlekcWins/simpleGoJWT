package dbModel

import (
	"golang.org/x/crypto/bcrypt"
	"simpleGoJWT/model"
)

type User struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	Email        string `json:"email" bson:"email"`
	Username     string `json:"username" bson:"username"`
	PasswordHash []byte `json:"-" bson:"password"`
}

func MapUserDto(u model.UserDTO) User {
	password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	return User{
		Email:        u.Email,
		Username:     u.Username,
		PasswordHash: password,
	}

}
