package mongo

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func New(ctx context.Context, cfg Config) (_ Wrapper, err error) {
	reqCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var clientOptions = cfg.getClientOpts()
	var client *mongo.Client
	if client, err = mongo.Connect(reqCtx, clientOptions); err != nil {
		return nil, fmt.Errorf("failed to create client to mongodb due to error %w", err)
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, fmt.Errorf("failed to create client to mongodb due to error %w", err)
	}

	var db = client.Database(cfg.Database)

	return &wrapper{
		db: db,
	}, nil
}

type (
	Wrapper interface {
		Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection
		Close() error
	}

	wrapper struct {
		db *mongo.Database
	}
)

func (w *wrapper) Close() error {
	return w.db.Client().Disconnect(context.Background())
}

func (w *wrapper) Collection(name string, opts ...*options.CollectionOptions) *mongo.Collection {
	return w.db.Collection(name, opts...)
}
