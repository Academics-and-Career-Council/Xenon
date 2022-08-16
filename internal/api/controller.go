package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/AnC-IITK/Xenon/internal/database"
	"github.com/AnC-IITK/Xenon/internal/gql"
	"github.com/AnC-IITK/Xenon/internal/services"
	"github.com/gocarina/gocsv"
	"github.com/gofiber/fiber/v2"
	acl "github.com/ory/keto/proto/ory/keto/acl/v1alpha1"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func whoami(c *fiber.Ctx) error {
	email := c.Get("X-Email")
	u := database.User{}
	err := services.BadgerDB.Get("user:"+email, u)
	if err == nil {
		return c.Status(200).JSON(u)
	}
	u, err = database.MongoClient.GetUser(email)
	if err != nil {
		return fiber.NewError(404, fmt.Sprint(err))
	}
	services.BadgerDB.Save("user:"+email, u)
	return c.Status(200).JSON(u)
}

func isGQLAllowed(c *fiber.Ctx) error {
	// Get Email Header Required for Authentication
	email := c.Get("X-Email")

	// Parse the request
	b := new(gql.GqlBody)
	err := json.Unmarshal(c.Request().Body(), b)
	if err != nil {
		// Return Unauthorized on Malformed Request
		log.Print("Malformed GraphQL Query")
		return c.SendStatus(403)
	}

	// Parse the GraphQL Request into a AST
	q, a := gql.Introspect(*b)

	// Fetch the corresponding Access Control List
	ketoACL := gql.ACL[q]

	// Execute the template to determine the requested resource
	tpl := new(bytes.Buffer)
	err = ketoACL.Object.Execute(tpl, a)
	if err != nil {
		log.Print("Malformed Golang Template")
		return c.SendStatus(403)
	}
	result := tpl.String()
	log.Println(ketoACL, email)
	// Check Permission using goRPC on Ory Keto
	log.Print(ketoACL.Namespace, result, ketoACL.Relation, email)
	allowed, err := services.CheckPermission(ketoACL.Namespace, result, ketoACL.Relation, email)
	log.Println(allowed, err)
	if err != nil {
		log.Print("Given Action/Subject/Resource does not exist!")
		return c.SendStatus(403)
	}
	if allowed {
		return c.SendStatus(200)
	}
	return c.SendStatus(403)
}

func Register(ctx *fiber.Ctx) error {
	u := strings.ReplaceAll(ctx.FormValue("username"), " ", "")
	r := ctx.FormValue("token")
	err := services.VerifyCaptcha(r)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to verify captcha")
	}
	code, err := services.KratosClient.CreateUser(u)
	if err != nil {
		return fiber.NewError(code, fmt.Sprint(err))
	}
	return ctx.Status(201).JSON(fiber.Map{"message": "Registered Successfully"})
}

func Recover(ctx *fiber.Ctx) error {
	u := strings.ReplaceAll(ctx.FormValue("username"), " ", "")
	log.Println(u)
	err := services.KratosClient.RecoverUser(u)
	if err != nil {
		return fiber.NewError(401, fmt.Sprint(err))
	}
	return ctx.SendStatus(200)
}

func UpdateAdministrators(ctx *fiber.Ctx) error {
	k := strings.ReplaceAll(ctx.FormValue("key"), " ", "")
	if k != viper.GetString("key") {
		return ctx.SendStatus(403)
	}
	f, err := ctx.FormFile("csv")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}
	b, err := f.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}
	reader := csv.NewReader(b)

	_, err = reader.Read()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}

	rt := []*acl.RelationTupleDelta{}
	mr := []mongo.WriteModel{}
	for i := 0; ; i = i + 1 {
		record, err := reader.Read()
		if err == io.EOF {
			break // reached end of the file
		} else if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
		} else if len(record) != 2 {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
		}

		t := &acl.RelationTuple{Namespace: "groups", Object: record[1], Relation: "member", Subject: &acl.Subject{Ref: &acl.Subject_Id{Id: record[0]}}}
		rt = append(rt, &acl.RelationTupleDelta{RelationTuple: t, Action: acl.RelationTupleDelta_INSERT})
		filter := bson.D{{Key: "username", Value: record[0]}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: record[1]}}}}
		mr = append(mr, &mongo.UpdateOneModel{Filter: filter, Update: update})
	}
	database.MongoClient.BulkWriteInStudents(mr)
	log.Println(mr)
	err = services.InsertTuples(ctx.Context(), rt)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, fmt.Sprint(err))
	}
	return ctx.SendStatus(200)
}

func InsertUsers(ctx *fiber.Ctx) error {
	k := strings.ReplaceAll(ctx.FormValue("key"), " ", "")
	if k != viper.GetString("key") {
		return ctx.SendStatus(403)
	}
	f, err := ctx.FormFile("csv")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}
	b, err := f.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}
	users := []*database.User{}
	if err := gocsv.Unmarshal(b, &users); err != nil { // Load clients from file
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprint(err))
	}

	rt := []*acl.RelationTupleDelta{}
	mr := []interface{}{}
	for _, user := range users {
		if !strings.Contains(user.EmailID, "@iitk.ac.in") || (user.Username+"@iitk.ac.in" != user.EmailID) {
			log.Println(fmt.Sprintf("Skipping %s, invalid parameters", user.Name))
		}
		t := &acl.RelationTuple{Namespace: "groups", Object: user.Role, Relation: "member", Subject: &acl.Subject{Ref: &acl.Subject_Id{Id: user.Username}}}
		rt = append(rt, &acl.RelationTupleDelta{RelationTuple: t, Action: acl.RelationTupleDelta_INSERT})
		user.Banned = false
		mr = append(mr, &user)
	}
	_, err = database.MongoClient.Users.Collection(viper.GetString("mongo.collection")).InsertMany(ctx.Context(), mr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, fmt.Sprint(err))
	}
	err = services.InsertTuples(ctx.Context(), rt)
	if err != nil {
		return fiber.NewError(fiber.StatusServiceUnavailable, fmt.Sprint(err))
	}
	return ctx.SendStatus(200)
}
