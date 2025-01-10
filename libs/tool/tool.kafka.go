package tool

import (
	"strings"

	"github.com/IBM/sarama"
)

// tools.SendKafkaProducerMessage("broker", "topic", "sync", "test")
// tools.HandlerKafkaConsumerMessage("broker", "topic")

// kafka producer
type KafkaProducer struct {
	Client map[string]sarama.AsyncProducer
}

func NewKafkaProducer(data map[string]string) *KafkaProducer {
	client := make(map[string]sarama.AsyncProducer, 0)

	for k, v := range data {
		addr := strings.Split(v, ",")
		client[k] = createKafkaProducerClient(addr)
	}

	return &KafkaProducer{
		Client: client,
	}
}

func createKafkaProducerClient(kfkConf []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = false
	config.Producer.Return.Successes = false

	producer, err := sarama.NewAsyncProducer(kfkConf, config)

	if err != nil {
		panic(err)
	}

	return producer
}

func (k *KafkaProducer) Close() {
	for _, v := range k.Client {
		v.Close()
	}
}

// kafka consumer
type KafkaConsumer struct {
	Client map[string]sarama.Consumer
}

func NewKafkaConsumer(data map[string]string) *KafkaConsumer {
	client := make(map[string]sarama.Consumer, 0)

	for k, v := range data {
		addr := strings.Split(v, ",")
		client[k] = createKafkaConsumerClient(addr)
	}

	return &KafkaConsumer{
		Client: client,
	}
}

func createKafkaConsumerClient(kfkConf []string) sarama.Consumer {
	consumer, err := sarama.NewConsumer(kfkConf, nil)

	if err != nil {
		panic(err)
	}

	return consumer
}

func (k *KafkaConsumer) Close() {
	for _, v := range k.Client {
		v.Close()
	}
}
