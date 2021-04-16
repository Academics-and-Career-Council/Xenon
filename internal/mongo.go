package internal

import (
	"context"
	"log"
	"strings"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	Users *mongo.Database
}

var MongoClient = &mongoClient{}

func ConnectMongo() {
	MongoClient.Users = connect(viper.GetString("mongo_url"), "primarydb")
}

func connect(url string, dbname string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Unable to Connect to MongoDB %v", err)
	}
	log.Printf("Connected to MongoDB! URL : %s", url)
	database := client.Database(dbname)
	return database
}

func (m mongoClient) getUser(email string, dest interface{}) error {
	e := strings.Split(email, "@")[0]
	filter := bson.D{{Key: "username", Value: e}}
	err := MongoClient.Users.Collection("students").FindOne(context.TODO(), filter).Decode(dest)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return err
}

func (m mongoClient) SetID(key string, id string, username string) error {
	u := &User{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: id}}}}
	err := m.Users.Collection("students").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return err
}

func (m mongoClient) CanRegister(username string) (bool, error) {
	u := &User{}
	name := strings.Replace(username, " ", "", -1)
	filter := bson.M{"username": name}
	err := m.Users.Collection("students").FindOne(context.TODO(), filter).Decode(u)
	// TODO Implement banning
	if err != nil {
		log.Printf("Unable to check access for %s: %v", name, err)
	}
	return true, err
}
