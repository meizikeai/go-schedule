package tool

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"

	"go-schedule/libs/types"

	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func HandleCron(time string, fn func()) *cron.Cron {
	res := cron.New()

	_, err := res.AddFunc(time, fn)

	if err != nil {
		log.Error(err)
	}

	res.Start()

	return res
}

func MarshalJson(date interface{}) []byte {
	res, err := json.Marshal(date)

	if err != nil {
		log.Error(err)
	}

	return res
}

func UnmarshalJson(date string) types.MapStringInterface {
	var res types.MapStringInterface

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func GetMODE() string {
	res := os.Getenv("GO_MODE")

	if res != "release" {
		res = "test"
	}

	return res
}
