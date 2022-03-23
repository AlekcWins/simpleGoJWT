package controllers

import (
	"github.com/gofiber/fiber/v2"
	"simpleGoJWT/database"
	"simpleGoJWT/database/dbModel"
	"simpleGoJWT/model"
	"simpleGoJWT/services"
)

func Register(c *fiber.Ctx) error {
	var userDTO model.UserDTO
	err := c.BodyParser(&userDTO)

	if err != nil {
		return err
	}

	user := dbModel.MapUserDto(userDTO)
	dbServices := services.GetDbServices(database.DB)
	id, err := dbServices.UserStorage.Create(user)
	if err != nil {
		_ = c.JSON(map[string]string{
			"message": "error create",
		})
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if id == "" {
		_ = c.JSON(map[string]string{
			"message": "this email already register",
		})
		return c.SendStatus(fiber.StatusBadRequest)
	}
	user.ID = id
	return c.JSON(user)
}
