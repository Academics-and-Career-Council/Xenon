package internal

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	kratos "github.com/ory/kratos-client-go/client"
	"github.com/ory/kratos-client-go/client/admin"
	kratosmodels "github.com/ory/kratos-client-go/models"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type Ory struct {
	*kratos.OryKratos
}

var OryClient = &Ory{}

func (ory *Ory) Connect() {
	ory.OryKratos = kratos.NewHTTPClientWithConfig(nil, &kratos.TransportConfig{
		Host:     viper.GetString("kratos.admin_url"),
		BasePath: "/",
		Schemes:  []string{"http"},
	})
}

func (ory *Ory) CreateUser(username string) error {
	ok, err := MongoClient.CanRegister(username)
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	if !ok {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	i, err := ory.OryKratos.Admin.CreateIdentity(&admin.CreateIdentityParams{Body: &kratosmodels.CreateIdentity{Traits: fiber.Map{"email": username + "@iitk.ac.in"}}, Context: context.TODO()})
	if err != nil {
		log.Printf("Unable to check regsitration : %v %v", err, err.Error())
		return err
	}
	id := i.GetPayload().ID
	fmt.Print(id)
	link, err := ory.OryKratos.Admin.CreateRecoveryLink(&admin.CreateRecoveryLinkParams{Body: &kratosmodels.CreateRecoveryLink{IdentityID: id}, Context: context.TODO()})
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	url := *link.Payload.RecoveryLink
	message := "Dear User,\nPlease use the following link to recover your account:\n" + url
	print(message)
	MongoClient.SetID("kid", string(*id), username)
	SendMail("New Registration", message, []string{username + "@iitk.ac.in"})
	return err
}

func (ory *Ory) RecoverUser(username string) error {
	filter := bson.D{{Key: "username", Value: username}}
	u := &User{}
	err := MongoClient.Users.Collection("ug").FindOne(context.TODO(), filter).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	id := kratosmodels.UUID(u.KratosID)
	fmt.Print(id)
	link, err := ory.OryKratos.Admin.CreateRecoveryLink(&admin.CreateRecoveryLinkParams{Body: &kratosmodels.CreateRecoveryLink{IdentityID: &id}, Context: context.TODO()})
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	url := *link.Payload.RecoveryLink
	message := "Dear User,\nPlease use the following link to recover your account:\n" + url
	print(message)

	SendMail("Account Recovery", message, []string{username + "@iitk.ac.in"})
	return err
}

// TODO
func (ory *Ory) DeleteUser(id string) error {
	_, err := ory.OryKratos.Admin.DeleteIdentity(&admin.DeleteIdentityParams{ID: id})
	return err
}
