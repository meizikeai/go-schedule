package tool

import (
	"fmt"
	"strings"

	"go-schedule/config"

	"github.com/IBM/sarama"
)

var fullProducerKafka map[string]sarama.AsyncProducer
var fullConsumerKafka map[string]sarama.Consumer

// producer
type KafkaProducer struct{}

func NewKafkaProducer() *KafkaProducer {
	return &KafkaProducer{}
}

func (t *Tools) HandleKafkaProducerClient() {
	config := config.GetKafkaConfig()
	result := make(map[string]sarama.AsyncProducer, len(config))

	for k, v := range config {
		addr := strings.Split(v, ",")
		result[k] = createKafkaProducerClient(addr)
	}

	fullProducerKafka = result

	t.Stdout("Kafka is Connected")
}

func createKafkaProducerClient(kfkConf []string) sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Return.Errors = false
	config.Producer.Return.Successes = false

	producer, err := sarama.NewAsyncProducer(kfkConf, config)

	if err != nil {
		panic(err.Error())
	}

	return producer
}

func (kp *KafkaProducer) GetKafkaProducerClient(key string) sarama.AsyncProducer {
	result := fullProducerKafka[key]
	return result
}

// demo
func (kp *KafkaProducer) SendKafkaProducerMessage(broker, topic, key, data string) {
	producer := kp.GetKafkaProducerClient(broker)

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(data),
	}

	producer.Input() <- message
}

// consumer
type KafkaConsumer struct{}

func NewKafkaConsumer() *KafkaConsumer {
	return &KafkaConsumer{}
}

func (t *Tools) HandleKafkaConsumerClient() {
	config := config.GetKafkaConfig()
	result := make(map[string]sarama.Consumer, len(config))

	for k, v := range config {
		addr := strings.Split(v, ",")
		result[k] = createKafkaConsumerClient(addr)
	}

	fullConsumerKafka = result
}

func createKafkaConsumerClient(kfkConf []string) sarama.Consumer {
	consumer, err := sarama.NewConsumer(kfkConf, nil)

	if err != nil {
		panic(err.Error())
	}

	return consumer
}

func (kc *KafkaConsumer) GetKafkaConsumerClient(key string) sarama.Consumer {
	result := fullConsumerKafka[key]
	return result
}

// demo
func (kc *KafkaConsumer) HandlerKafkaConsumerMessage(broker, topic string) {
	consumer := kc.GetKafkaConsumerClient(broker)
	partitionList, err := consumer.Partitions(topic)

	if err != nil {
		panic(err.Error())
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)

		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)
	}
}

func (t *Tools) CloseKafka() {
	for _, v := range fullProducerKafka {
		v.Close()
	}

	for _, v := range fullConsumerKafka {
		v.Close()
	}

	t.Stdout("Kafka is Close")
}
