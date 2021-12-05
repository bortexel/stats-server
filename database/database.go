package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client   *mongo.Client
	Database *mongo.Database
)

func InitDatabase(uri string) (err error) {
	Client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return
	}

	Database = Client.Database("stats")
	return
}
