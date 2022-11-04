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
	clientOptions := options.Client().ApplyURI(uri)
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	Database = Client.Database("stats")
	return
}
