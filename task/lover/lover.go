package lover

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-schedule/models"

	log "github.com/sirupsen/logrus"
)

func HandleLoverGift() {
	var sum int64 = 0
	start, ch := time.Now(), make(chan int64)

	people := make(map[int]int)

	day := time.Now().Format("20060102")
	day = "20200205"

	back, _ := models.GetLoverGift(day)

	for _, val := range back {
		if people[val.Uid] > 0 {
			beans, _ := strconv.Atoi(val.Beans)
			value := people[val.Uid] + beans
			people[val.Uid] = value
		} else {
			beans, _ := strconv.Atoi(val.Beans)
			people[val.Uid] = beans
		}
	}

	log.Info(people)

	result := make([]string, 0)

	if len(people) > 0 {
		for k, v := range people {
			key := strconv.Itoa(k)
			value := strconv.Itoa(v)
			str := []string{"(", key, ",", value, ",", day, ")"}

			result = append(result, strings.Join(str, ""))
		}
	}
	// log.Info(result)

	if len(result) > 0 {
		for _, v := range result {
			go func(value string) {
				id, _ := models.UpdateLoverData("lover_daily_beans", value, "beans = VALUES(beans)")
				ch <- id
			}(v)
		}

		sum += <-ch

		end := time.Now()
		log.Info(fmt.Sprintf("result is %d, time is %s", sum, end.Sub(start)))
	}
}
