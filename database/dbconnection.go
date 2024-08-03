package database

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	errDotenverr = errors.New("error while loading env file")
)

func DbInstance() (*mongo.Client, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, errDotenverr
	}
	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, err
}

func OpenCollection(client *mongo.Client, collectionname string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("cluster0").Collection(collectionname)
	return collection
}
