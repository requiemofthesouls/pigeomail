package database

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("data not found")

func NewMongoDBClient(cfg *Config) (client *mongo.Client, err error) {
	var uri = fmt.Sprintf("mongodb://%s:%d/?connect=direct", cfg.Hostname, cfg.Port)
	var creds = options.Credential{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	var clientOpts = options.Client().ApplyURI(uri).SetAuth(creds)

	if client, err = mongo.Connect(context.TODO(), clientOpts); err != nil {
		return nil, err
	}

	return client, nil
}
