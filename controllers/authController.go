package controllers

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"simpleGoJWT/database"
	"simpleGoJWT/database/dbModel"
	"simpleGoJWT/model"
	"simpleGoJWT/services"
)

var dbServices = services.GetDbServices(database.DB)

func Login(c *fiber.Ctx) error {
	loginDTO := model.LoginDTO{}

	err := c.BodyParser(&loginDTO)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "invalid login data",
		})
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}

	user, err := dbServices.UserStorage.FindOneWithEmail(loginDTO.Email)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "user for this email not found",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(loginDTO.Password)); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	jwt := services.Jwt{}
	token, err := jwt.CreateToken(user.ID)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "unable to create access token",
		})

		return c.SendStatus(fiber.StatusInternalServerError)
	}
	saveToken(c, user.ID, token)
	return c.JSON(token)

}

func RefreshToken(c *fiber.Ctx) error {
	token := model.TokenDTO{}
	err := c.BodyParser(&token)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "invalid token",
		})

		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	jwt := services.Jwt{}
	userId, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "invalid token",
		})

		return c.SendStatus(fiber.StatusBadRequest)
	}

	tokenFound, err := dbServices.TokenStorage.FindOneWithUserID(userId)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "refresh token not found",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if tokenFound.IsBlockToken {
		_ = c.JSON(map[string]string{
			"message": "this refresh token blocked please login again",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !tokenFound.CompareTokenHash(token.RefreshToken) {
		_ = c.JSON(map[string]string{
			"message": "this refresh token does not match",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err = jwt.CreateToken(userId)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "unable to create access token",
		})

		return c.SendStatus(fiber.StatusInternalServerError)
	}
	saveToken(c, userId, token)
	return c.JSON(token)

}

func BlockRefreshToken(c *fiber.Ctx) error {
	var data map[string]string
	err := c.BodyParser(&data)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "invalid email data",
		})
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	value, ok := data["email"]
	if !ok {
		_ = c.JSON(map[string]string{
			"message": "invalid email data",
		})
		return c.SendStatus(fiber.StatusUnprocessableEntity)
	}
	user, err := dbServices.UserStorage.FindOneWithEmail(value)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "not found user with this email",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	tokenFound, err := dbServices.TokenStorage.FindOneWithUserID(user.ID)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "refresh token not found",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	tokenFound.IsBlockToken = true
	_, err = dbServices.TokenStorage.CreateOrUpdate(tokenFound)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "refresh token cant be saves mongo",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(map[string]string{
		"message": "success",
	})
}

func saveToken(c *fiber.Ctx, userId string, token model.TokenDTO) {
	t := dbModel.NewToken(userId, token.RefreshToken)
	_, err := dbServices.TokenStorage.CreateOrUpdate(*t)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "failed save access token to mongo",
		})
		_ = c.SendStatus(fiber.StatusInternalServerError)
	}
}
