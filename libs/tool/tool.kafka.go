package tool

import (
	"fmt"
	"log"
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

// Consumer represents a Sarama consumer group consumer
type KafkaConsumerGroup struct {
	Ready    chan bool
	Callback func([]byte) error
}

func NewKafkaConsumerGroupConsumer(callback func([]byte) error) *KafkaConsumerGroup {
	return &KafkaConsumerGroup{
		Ready:    make(chan bool),
		Callback: callback,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *KafkaConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *KafkaConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *KafkaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				log.Panicf("message channel was closed")
				return nil
			}
			// fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", message.Partition, message.Offset, message.Key, string(message.Value))

			err := c.Callback(message.Value)

			if err == nil {
				session.MarkMessage(message, "")
			}

		// Should return when `session.Context()` is done.
		// If not, will raise `ErrRebalanceInProgress` or `read tcp <ip>:<port>: i/o timeout` when kafka rebalance. see:
		// https://github.com/IBM/sarama/issues/1192
		case <-session.Context().Done():
			return nil
		}
	}
}
