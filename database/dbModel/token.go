package dbModel

import (
	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	ID           string `json:"id" bson:"_id,omitempty"`
	UserUId      string `json:"user_u_id" bson:"user_id"`
	Token        string `json:"-" bson:"token"`
	IsBlockToken bool   `json:"is_block_token" bson:"is_block_token"`
}

func (t *Token) SetToken(token string) {
	tokenHash, _ := bcrypt.GenerateFromPassword([]byte(token), 14)
	t.Token = string(tokenHash)
}

func (t *Token) CompareTokenHash(token string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(token)); err != nil {
		return false
	}
	return true
}

func NewToken(useId, token string) *Token {
	tokenCreate := &Token{
		UserUId:      useId,
		IsBlockToken: false,
	}
	tokenCreate.SetToken(token)

	return tokenCreate

}
