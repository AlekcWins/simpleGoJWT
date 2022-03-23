package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"io"
	"simpleGoJWT/config"
	"simpleGoJWT/model"
	"time"
)

var secretKey string

func init() {
	secretKey = config.GetConfig().JWT.SecretKeyJWT
}

type Jwt struct {
}

func (j Jwt) CreateToken(userid string) (model.TokenDTO, error) {
	var err error
	claims := jwt.MapClaims{}
	claims["user_id"] = userid
	claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	jwtToken := model.TokenDTO{}

	jwtToken.AccessToken, err = token.SignedString([]byte(secretKey))
	if err != nil {
		return jwtToken, err
	}

	return j.createRefreshToken(jwtToken)
}

func (Jwt) ValidateToken(accessToken string) (id string, err error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secretKey), nil
	})
	if err != nil {
		return id, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id = payload["user_id"].(string)

		return id, nil
	}

	return id, errors.New("invalid token")
}

func (Jwt) ValidateRefreshToken(model model.TokenDTO) (id string, err error) {
	sha := sha1.New()
	_, err = io.WriteString(sha, secretKey)
	if err != nil {
		return id, errors.New("error write secret key")
	}

	salt := string(sha.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		return id, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return id, err
	}

	data, err := base64.URLEncoding.DecodeString(model.RefreshToken)
	if err != nil {
		return id, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return id, err
	}

	if string(plain) != model.AccessToken {
		return id, errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err := parser.ParseUnverified(model.AccessToken, claims)

	if err != nil {
		return id, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return id, errors.New("invalid token")
	}

	id = payload["user_id"].(string)
	return id, nil
}

func (Jwt) createRefreshToken(token model.TokenDTO) (t model.TokenDTO, err error) {
	sha := sha1.New()
	_, err = io.WriteString(sha, secretKey)
	if err != nil {
		return t, errors.New("error write secret key")
	}

	salt := string(sha.Sum(nil))[0:16]
	block, err := aes.NewCipher([]byte(salt))
	if err != nil {
		fmt.Println(err.Error())

		return token, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return token, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return token, err
	}

	token.RefreshToken = base64.URLEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(token.AccessToken), nil))

	return token, nil
}
