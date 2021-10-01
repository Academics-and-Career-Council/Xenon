package cmd

import (
	"github.com/AnC-IITK/Xenon/internal/api"
	"github.com/AnC-IITK/Xenon/internal/database"
	"github.com/AnC-IITK/Xenon/internal/gql"
	"github.com/AnC-IITK/Xenon/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Fiber Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		gql.InitializeACL()
		database.ConnectMongo()
		services.ConntectKratos()
		err := Serve()
		return err
	},
}

func Serve() error {
	app := fiber.New()
	api.SetupRoutes(app)
	err := app.Listen(":5000")
	return err
}
