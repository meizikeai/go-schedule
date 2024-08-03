package models

import (
	"context"

	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

func GetUserName() (result string, err error) {
	pool := tools.GetRedisClient("users.master")
	back, err := pool.HGet(ctx, "u:644", "name").Result()
	// back, err := pool.HGetAll(ctx, "u:644").Result()

	if err != nil {
		log.Error(err)
	}

	result = back

	return result, err
}
