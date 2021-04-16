package main

import (
	"log"

	"github.com/AnC-IITK/Xenon/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.ReadInConfig()
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
