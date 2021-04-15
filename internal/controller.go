package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mikespook/gorbac"
	"github.com/spf13/viper"
)

func whoami(c *fiber.Ctx) error {
	email := c.Get("X-IITK-Email")
	u, err := MongoClient.getUser(email)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"Role": "USER"})
	}
	return c.Status(200).JSON(fiber.Map{"Role": u.Role})
}

func isGQLAllowed(c *fiber.Ctx) error {
	email := c.Get("X-IITK-Email")
	u, _ := MongoClient.getUser(email)
	b := new(gqlBody)
	json.Unmarshal(c.Request().Body(), b)
	if PermissionManager.IsGranted(u.Role, gorbac.NewStdPermission(b.Query), nil) && !u.Banned {
		if b.Query == "getReviews" && u.Credits < viper.GetInt("credits_unlock") {
			vars := map[string]string{}
			json.Unmarshal([]byte(b.Variables), &vars)
			res, _ := IntrospectGetReviews(email, vars["course"])
			if res {
				return c.SendStatus(200)
			} else {
				return c.SendStatus(403)
			}
		} else {
			return c.SendStatus(200)
		}
	} else {
		return c.SendStatus(403)
	}
}

func Register(ctx *fiber.Ctx) error {
	u := strings.ReplaceAll(ctx.FormValue("username"), " ", "")
	log.Println(u)
	err := OryClient.CreateUser(u)
	if err != nil {
		return fiber.NewError(401, fmt.Sprint(err))
	}
	return ctx.SendStatus(201)
}

func Recover(ctx *fiber.Ctx) error {
	u := strings.ReplaceAll(ctx.FormValue("username"), " ", "")
	log.Println(u)
	err := OryClient.RecoverUser(u)
	if err != nil {
		return fiber.NewError(401, fmt.Sprint(err))
	}
	return ctx.SendStatus(200)
}
