package tool

import (
	"context"
	"time"

	"go-schedule/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var fullDbMongoDB map[string][]*mongo.Client

func GetMongoClient() *mongo.Client {
	mongo := fullDbMongoDB["default.master"]
	count := GetRandmod(len(mongo))

	result := mongo[count]

	return result
}

func GetMongoCollection(database, table string) *mongo.Collection {
	client := GetMongoClient()

	db := client.Database(database)
	result := db.Collection(table)

	return result
}

func HandleMongodbClient() {
	client := make(map[string][]*mongo.Client)

	local := config.GetMongodbConfig()

	for key, val := range local {
		m := key + ".master"
		for _, v := range val {
			clients := createMongodbClient(v)
			client[m] = append(client[m], clients)
		}
	}

	fullDbMongoDB = client
}

func createMongodbClient(uri string) *mongo.Client {
	config := options.Client().ApplyURI(uri)
	ctx, cancel := getMongoConfig()

	client, err := mongo.Connect(ctx, config)

	if err != nil {
		panic(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		panic(err)
	}

	defer cancel()

	return client
}

func getMongoConfig() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.TODO(), 4*time.Second)
	return ctx, cancel
}

func CloseMongoDB() {
	for _, val := range fullDbMongoDB {
		for _, v := range val {
			err := v.Disconnect(context.TODO())

			if err != nil {
				// panic(err)
			}
		}
	}

	Stdout("MongoDB is Close")
}
