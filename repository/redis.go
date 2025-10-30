package repository

import (
	"encoding/json"
	"fmt"

	"go-schedule/libs/types"
)

type ModelsRedis struct{}

func NewModelsRedis() *ModelsRedis {
	return &ModelsRedis{}
}

func (m *ModelsRedis) SaveBinlogPosition(name, data string) {
	db := tools.GetRedisClient("default")

	key := fmt.Sprintf("clio:sync:%s", name)
	_, err := db.Set(ctx, key, data, 0).Result()

	if err != nil {
		fmt.Println(err.Error())
	}
}

func (m *ModelsRedis) GetBinlogPosition(name string) types.BinlogPosition {
	db := tools.GetRedisClient("default")

	result := types.BinlogPosition{}

	key := fmt.Sprintf("clio:sync:%s", name)
	data, err := db.Get(ctx, key).Result()

	if err != nil {
		return result
	}

	err = json.Unmarshal([]byte(data), &result)

	if err != nil {
		return result
	}

	return result
}
