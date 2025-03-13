package tool

import (
	"fmt"
	"strings"

	"go-schedule/config"

	"github.com/IBM/sarama"
)

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

// kafka consumer group
type KafkaConsumerCroup struct {
	Client map[string]sarama.ConsumerGroup
}

func NewKafkaConsumerGroup(data map[string]string) *KafkaConsumerCroup {
	client := make(map[string]sarama.ConsumerGroup, 0)

	for k, v := range data {
		addr := strings.Split(v, ",")
		option := config.GetKafkaConsumerGroupConfig(k)
		client[k] = createKafkaConsumerGroupClient(addr, option)
	}

	return &KafkaConsumerCroup{
		Client: client,
	}
}
func createKafkaConsumerGroupClient(kfkConf []string, options map[string]string) sarama.ConsumerGroup {
	config := sarama.NewConfig()

	if options["version"] == "" {
		options["version"] = sarama.DefaultVersion.String()
	}
	version, err := sarama.ParseKafkaVersion(options["version"])

	if err != nil {
		panic(fmt.Sprintf("Error parsing Kafka version: %v", err))
	}
	config.Version = version

	switch options["assignor"] {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "roundrobin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		panic(fmt.Sprintf("Unrecognized consumer group partition assignor: %s", options["assignor"]))
	}

	if options["oldest"] == "true" {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	consumer, err := sarama.NewConsumerGroup(kfkConf, options["groupID"], config)

	if err != nil {
		panic(err)
	}

	return consumer
}

func (k *KafkaConsumerCroup) Close() {
	for _, v := range k.Client {
		v.Close()
	}
}
