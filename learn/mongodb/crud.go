package mongodb

import (
	"context"
	"time"

	"github.com/kuangcp/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	logger.Info(client, err)
	ctx, caf := context.WithTimeout(context.Background(), 20*time.Second)
	caf()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		logger.Error(err)
	}
	database := client.Database("test-data")
	database.CreateCollection(ctx, "xxx")
}
