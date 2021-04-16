package cmd

import (
	"github.com/AnC-IITK/Xenon/internal"
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
		internal.PermissionsInit()
		internal.ConnectMongo()
		internal.OryClient.Connect()
		internal.Init()
		err := Serve()
		return err
	},
}

func Serve() error {
	app := fiber.New()
	internal.SetupRoutes(app)
	err := app.Listen(":5000")
	return err
}
