package test

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"go-schedule/libs/tool"
	"go-schedule/models"

	log "github.com/sirupsen/logrus"
)

var job int = 0

func OneJob() {
	job++

	log.Info(job, ", cron is running...")
}

func TwoJob() {
	rand.Seed(time.Now().Unix())

	count := rand.Intn(1000000)

	person := []string{"guest" + strconv.Itoa(count) + "@test.com", "guest", "汉族", "男", "11010199812187756", "13412345678", "北京市朝阳区百子湾路苹果社区B区", "100000"}
	lastId, _ := models.AddPerson(person)

	log.Info(strconv.Itoa(count), lastId)
}

func ThreeJob() {
	var wg = sync.WaitGroup{}

	wg.Add(5)

	for i := 0; i < 5; i++ {
		go func(s int) {
			log.Info(s)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func FourJob() {
	data := make(chan string)
	exit := make(chan bool)

	go func() {
		for d := range data {
			log.Info(d)
		}

		exit <- true
	}()

	data <- "a"
	data <- "b"
	data <- "c"

	close(data)

	<-exit
}

func FiveJob() {
	a := make(chan models.Person)
	b := make(chan models.Person)

	go func() {
		one, _ := models.GetPerson("admin@bank.com")

		a <- one
	}()

	go func() {
		two, _ := models.GetPerson("test@bank.com")

		b <- two
	}()

	log.Info(string(tool.MarshalJson(<-a)))
	log.Info(string(tool.MarshalJson(<-b)))

	close(a)
	close(b)
}

func SixJob() {
	data := [2]string{"admin@bank.com", "test@bank.com"}
	ch := make(chan models.Person, 2)

	for i := 0; i < len(data); i++ {
		go func(i int) {
			one, _ := models.GetPerson(data[i])
			ch <- one
		}(i)
	}

	for v := range ch {
		log.Info(string(tool.MarshalJson(v)))
	}

	close(ch)
}
