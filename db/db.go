package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection interface {
}

var Db *mongo.Collection
var DbCtx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(DbCtx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(DbCtx, nil)
	if err != nil {
		log.Fatal(err)
	}

	Db = client.Database("my-yarn-api").Collection("yarns")
}
