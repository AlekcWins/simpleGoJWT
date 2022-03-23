package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	log "github.com/sirupsen/logrus"
	"simpleGoJWT/config"
	"simpleGoJWT/routs"
)

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routs.Setup(app)

	app.Use(logger.New())

	cfg := config.GetConfig()
	listenConf := cfg.Listen

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", listenConf.Host, listenConf.Port)))

}
