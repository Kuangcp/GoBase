package mongodemo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func connect() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.16.8.3:27017"))
	log.Print(client, err)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	print(err)

}
