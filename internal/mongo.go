package internal

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	*mongo.Database
}

var MongoClient mongoClient

func (m mongoClient) Connect() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://anc:courses@cluster0.x6adj.mongodb.net/students?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		print(err.Error())
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		print(err.Error())
	}
	log.Printf("Connected to MongoDB!")
	database := client.Database("students")
	m.Database = database
}

func (m mongoClient) getUser(email string, dest interface{}) error {
	e := strings.Split(email, "@")[0]
	filter := bson.D{{Key: "username", Value: e}}
	err := MongoClient.Collection("ug").FindOne(context.TODO(), filter).Decode(dest)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return err
}

func (m mongoClient) SetID(key string, id string, username string) error {
	u := &User{}
	filter := bson.D{{Key: "username", Value: username}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: id}}}}
	err := m.Collection("ug").FindOneAndUpdate(context.TODO(), filter, update).Decode(u)
	if err != nil {
		log.Printf("Unable to check access : %v", err)
	}
	return err
}

func (m mongoClient) CanRegister(username string) (bool, error) {
	u := &User{}
	name := strings.Replace(username, " ", "", -1)
	filter := bson.M{"username": name}
	err := m.Collection("ug").FindOne(context.TODO(), filter).Decode(u)
	// TODO Implement banning
	if err != nil {
		log.Printf("Unable to check access for %s: %v", name, err)
	}
	return true, err
}
