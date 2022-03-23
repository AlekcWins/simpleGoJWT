package routs

import (
	"github.com/gofiber/fiber/v2"
	"simpleGoJWT/controllers"
)

func InitUserRoute(app fiber.Router) {
	app.Post("/register", controllers.Register)
}
