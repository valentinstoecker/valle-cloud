package db

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func init() {
	con_str, set := os.LookupEnv("MONGO_STRING")
	if !set {
		con_str = "mongodb://db:27017"
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(con_str))
	if err != nil {
		panic(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}
	DB = client.Database("db")
}
