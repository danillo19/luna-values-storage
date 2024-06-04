package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func NewMongoClient(ctx context.Context, url string) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(url)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	// Ping the primary
	//TODO: could be problems with Primary
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	fmt.Println("MONGODB: Successfully connected and pinged.")

	return client, err
}
