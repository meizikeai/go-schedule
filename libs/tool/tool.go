package tool

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/robfig/cron/v3"
)

type Tools struct{}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) HandleCron(time string, fn func()) *cron.Cron {
	res := cron.New()

	_, err := res.AddFunc(time, fn)

	if err != nil {
		fmt.Println(err)
	}

	res.Start()

	return res
}

func (t *Tools) GetApiHost(key string) string {
	result := ""

	for k, v := range zookeeperApi {
		if k == key {
			i := t.GetRandmod(len(v))
			result = v[i]
		}
	}

	return result
}

func (t *Tools) GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func (t *Tools) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
