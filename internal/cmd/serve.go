package cmd

import (
	"github.com/AnC-IITK/Xenon/internal/api"
	"github.com/AnC-IITK/Xenon/internal/database"
	"github.com/AnC-IITK/Xenon/internal/gql"
	"github.com/AnC-IITK/Xenon/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"log"
	zmq "github.com/pebbe/zmq4"
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
		services.ConnectKeto()
		services.Init()
		err := Serve()
		return err
	},
}

type Input struct {
    Text string `json:"text"`
}

func Serve() error {
	// Create a new fiber instance with custom config
	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// Send custom error page
			err = ctx.Status(code).JSON(fiber.Map{"message": err.Error()})
			if err != nil {
				// In case the SendFile fails
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal Server Error"})
			}

			// Return from handler
			return nil
		},
	})
	zctx, _ := zmq.NewContext()

    zm, _ := zctx.NewSocket(zmq.REQ)
    zm.Connect("tcp://localhost:5555")
	input := new(Input)
	log.Println("Gocron Working for ZMQ")
	zm.Send(input.Text, 0)

	if msg, err := zm.Recv(0); err != nil {
		panic(err)
	} else {
		log.Println("Sent message %s", msg)
	}
	api.SetupRoutes(app)
	err := app.Listen(":5010")
	return err
}
