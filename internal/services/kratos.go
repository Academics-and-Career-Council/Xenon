package services

import (
	"context"
	"log"

	"github.com/AnC-IITK/Xenon/internal/database"
	"github.com/gofiber/fiber/v2"
	client "github.com/ory/kratos-client-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
)

type kratos struct {
	*client.APIClient
}

var KratosClient kratos

func ConntectKratos() {
	KratosClient = kratos{client.NewAPIClient(&client.Configuration{
		Servers: client.ServerConfigurations{client.ServerConfiguration{URL: viper.GetString("kratos.admin_url")}},
	})}
}

func (k kratos) DeleteIdentity(id string) {
	k.V0alpha1Api.AdminDeleteIdentity(context.TODO(), id).Execute()
}

func (k kratos) RecoverUser(username string) error {
	filter := bson.D{{Key: "username", Value: username}}
	u := &database.User{}
	err := database.MongoClient.Users.Collection("ug").FindOne(context.TODO(), filter).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}

	link, _, err := k.V0alpha1Api.AdminCreateSelfServiceRecoveryLink(context.TODO()).AdminCreateSelfServiceRecoveryLinkBody(client.AdminCreateSelfServiceRecoveryLinkBody{IdentityId: u.KratosID}).Execute()
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return err
	}
	url := link.GetRecoveryLink()
	message := "<p>Dear User,<br/>Please use the following link to recover your account:<br/><a href=\"" + url + "\">" + url + "</a></p>"

	err = SendMail("Account Recovery", message, []string{username + "@iitk.ac.in"})
	return err
}

func (k kratos) CreateUser(username string) (int, error) {
	ok, err := database.MongoClient.CanRegister(username)
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return 500, err
	}
	if !ok {
		log.Printf("Unable to check regsitration : %v", err)
		return 404, err
	}
	b := k.V0alpha1Api.AdminCreateIdentity(context.TODO())
	i, r, err := b.AdminCreateIdentityBody(client.AdminCreateIdentityBody{Traits: fiber.Map{"email": username + "@iitk.ac.in"}}).Execute()
	if err != nil {
		log.Printf("Unable to check regsitration : %v %v", err, r)
		return r.StatusCode, err
	}
	database.MongoClient.SetID("kid", i.Id, username)
	link, r, err := k.V0alpha1Api.AdminCreateSelfServiceRecoveryLink(context.TODO()).AdminCreateSelfServiceRecoveryLinkBody(client.AdminCreateSelfServiceRecoveryLinkBody{IdentityId: i.Id}).Execute()
	if err != nil {
		log.Printf("Unable to check regsitration : %v", err)
		return r.StatusCode, err
	}
	url := link.GetRecoveryLink()
	message := "Dear User,\nPlease use the following link to recover your account:\n" + url
	err = SendMail("New Registration", message, []string{username + "@iitk.ac.in"})
	if err != nil {
		log.Printf("SMTP Error : %v", err)
		return 500, err
	}
	return 200, err
}
