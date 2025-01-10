package tool

import (
	"context"
	"time"

	"go-schedule/libs/types"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client map[string]*mongo.Client
}

func NewMongoDB(data map[string]types.ConfMongoDB) *MongoDB {
	client := make(map[string]*mongo.Client, 0)

	for k, v := range data {
		m := k + ".master"
		s := k + ".slave"

		master := createMongoDBClient(v.Master)
		client[m] = master

		slave := createMongoDBClient(v.Slave)
		client[s] = slave
	}

	return &MongoDB{
		Client: client,
	}
}

func createMongoDBClient(uri string) *mongo.Client {
	ctx, cancel := mongoConfig()
	defer cancel()

	config := options.Client().ApplyURI(uri)

	config.SetMaxPoolSize(300)
	config.SetMinPoolSize(150)

	client, err := mongo.Connect(ctx, config)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	return client
}

func mongoConfig() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	return ctx, cancel
}

func (t *MongoDB) Close() {
	for _, v := range t.Client {
		err := v.Disconnect(context.TODO())

		if err != nil {
			// panic(err)
		}
	}
}
