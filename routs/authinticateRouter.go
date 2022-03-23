package routs

import (
	"github.com/gofiber/fiber/v2"
	"simpleGoJWT/controllers"
)

func InitAuthenticationRoute(app fiber.Router) {
	app.Post("/refresh/", controllers.RefreshToken)
	app.Post("/login/", controllers.Login)
	app.Post("/block_token/", controllers.BlockRefreshToken)
}
