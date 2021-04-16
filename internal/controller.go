package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func whoami(c *fiber.Ctx) error {
	email := c.Get("X-IITK-Email")
	u, err := getUser(email)
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"Role": "USER"})
	}
	return c.Status(200).JSON(fiber.Map{"Role": u.Role})
}

func isGQLAllowed(c *fiber.Ctx) error {
	email := c.Get("X-IITK-Email")
	u, err := getUser(email)
	if err != nil {
		return c.SendStatus(403)
	}
	b := new(gqlBody)
	err = json.Unmarshal(c.Request().Body(), b)
	if err != nil {
		log.Print("Malformed GraphQL Query")
		return c.SendStatus(403)
	}
	q := Introspect(b.Query)
	if PermissionManager.rbac.IsGranted(u.Role, PermissionManager.permissions[q], nil) && !u.Banned {
		if q == "getReviews" && u.Credits < viper.GetInt("credits_unlock") {
			res := IntrospectGetReviews(*u, b.Variables["course"])
			if res {
				log.Printf("%s granted access for %s", email, q)
				return c.SendStatus(200)
			} else {
				log.Printf("%s blocked access for %s : Insufficient Credits", email, b.Variables["course"])
				return c.SendStatus(403)
			}
		} else {
			log.Printf("%s granted access for %s", email, q)
			return c.SendStatus(200)
		}
	} else {
		log.Printf("%s blocked access for %s : Permission Denied", email, q)
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
