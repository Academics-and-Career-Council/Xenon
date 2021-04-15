package main

import (
	"log"

	"github.com/AnC-IITK/Xenon/internal"
	"github.com/gofiber/fiber/v2"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	internal.PermissionManager.Init()
	internal.MongoClient.Connect()
	internal.OryClient.Connect()
	internal.BadgerDB.Init()
}

func main() {
	app := fiber.New()
	internal.SetupRoutes(app)
	app.Listen(":5000")
}
