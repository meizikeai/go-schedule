package tool

import (
	"context"
	"time"

	"go-schedule/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fullMongoDB map[string]*mongo.Client

func (t *Tools) GetMongoCollection(key, database, table string) *mongo.Collection {
	client := fullMongoDB[key]

	db := client.Database(database)
	collection := db.Collection(table)

	return collection
}

func (t *Tools) HandleMongoDBClient() {
	clients := make(map[string]*mongo.Client)

	local := config.GetMongodbConfig()

	for k, v := range local {
		m := k + ".master"
		s := k + ".slave"

		master := createMongoDBClient(v.Master)
		clients[m] = master

		slave := createMongoDBClient(v.Slave)
		clients[s] = slave
	}

	fullMongoDB = clients

	t.Stdout("MongoDB is Connected")
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

func (t *Tools) CloseMongoDB() {
	for _, v := range fullMongoDB {
		err := v.Disconnect(context.TODO())

		if err != nil {
			// panic(err)
		}
	}

	t.Stdout("MongoDB is Close")
}
