package main

import (
	"log"

	"github.com/AnC-IITK/Xenon/internal"
	"github.com/gofiber/fiber/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	internal.PermissionsInit()
	internal.ConnectMongo()
	internal.OryClient.Connect()
	internal.Init()
}

func main() {
	app := fiber.New()
	internal.SetupRoutes(app)
	app.Listen(":5000")
}
