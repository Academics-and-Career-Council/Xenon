package api

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/whoami", whoami)
	app.Post("/isAllowed/graphql", isGQLAllowed)
	app.Post("/isAllowed/rest", isRESTAllowed)
	app.Post("/register", Register)
	app.Post("/recover", Recover)
	app.Post("/admin/update", UpdateAdministrators)
	app.Put("/admin/insert", InsertUsers)
}
