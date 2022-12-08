package tool

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
)

// var kafkaConfig = map[string]string{
// 	"broker": "127.0.0.1:9092,127.0.0.1:9092,127.0.0.1:9092",
// }

var fullProducerKafka map[string]sarama.AsyncProducer
var fullConsumerKafka map[string]sarama.Consumer

// producer
func HandleKafkaProducerClient() {
	kafkaConfig := GetZookeeperStringConfig()
	result := make(map[string]sarama.AsyncProducer, len(kafkaConfig))

	for k, v := range kafkaConfig {
		addr := strings.Split(v, ",")
		result[k] = createKafkaProducerClient(addr)
	}

	fullProducerKafka = result
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

func GetKafkaProducerClient(key string) sarama.AsyncProducer {
	result := fullProducerKafka[key]
	return result
}

// demo
func SendKafkaProducerMessage(broker, topic, key, data string) {
	producer := GetKafkaProducerClient(broker)

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.StringEncoder(data),
	}

	producer.Input() <- message
}

// consumer
func HandleKafkaConsumerClient() {
	kafkaConfig := GetZookeeperStringConfig()
	result := make(map[string]sarama.Consumer, len(kafkaConfig))

	for k, v := range kafkaConfig {
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

func GetKafkaConsumerClient(key string) sarama.Consumer {
	result := fullConsumerKafka[key]
	return result
}

func CloseKafka() {
	for _, v := range fullProducerKafka {
		v.Close()
	}

	for _, v := range fullConsumerKafka {
		v.Close()
	}

	Stdout("Kafka Close")
}

// demo
func HandlerKafkaConsumerMessage(broker, topic string) {
	consumer := GetKafkaConsumerClient(broker)
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
