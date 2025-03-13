package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
)

type ModelsKafka struct{}

func NewModelsKafka() *ModelsKafka {
	return &ModelsKafka{}
}

func (m *ModelsKafka) SendKafkaProducerMessage(broker, topic, data string) {
	producer := tools.GetKafkaProducerClient(broker)

	if producer == nil {
		return
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(data),
	}

	producer.Input() <- message
}

func (m *ModelsKafka) ReadKafkaConsumerMessage(broker, topic string) {
	consumer := tools.GetKafkaConsumerClient(broker)
	partitionList, err := consumer.Partitions(topic)

	if err != nil {
		fmt.Println(err)
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)

		if err != nil {
			fmt.Printf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}

		go func(sarama.PartitionConsumer) {
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d Offset:%d Key:%v Value:%v\n", msg.Partition, msg.Offset, msg.Key, string(msg.Value))
			}
		}(pc)
	}
}

// https://github.com/IBM/sarama/blob/main/examples/consumergroup/main.go
func (m *ModelsKafka) ReadKafkaConsumerGroupMessage() {
	keepRunning := true
	tools.Stdout("Starting a new Sarama consumer")

	// Setup a new Sarama consumer group
	consumer := Consumer{
		ready: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	client := tools.GetKafkaConsumerGroupClient("default")

	consumptionIsPaused := false

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{"nyx_traffic"}, &consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	tools.Stdout("Sarama consumer up and running!...")

	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			tools.Stdout("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			tools.Stdout("terminating: via signal")
			keepRunning = false
		case <-sigusr1:
			toggleConsumptionFlow(client, &consumptionIsPaused)
		}
	}

	cancel()
	wg.Wait()

	if err := client.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		tools.Stdout("Resuming consumption")
	} else {
		client.PauseAll()
		tools.Stdout("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/IBM/sarama/blob/main/consumer_group.go#L27-L29
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				tools.Stdout("message channel was closed")
				return nil
			}
			fmt.Printf("Topic:%s Partition:%d Offset:%d Key:%v Value:%v\n", message.Topic, message.Partition, message.Offset, message.Key, string(message.Value))

			err := handleMessage(message.Value)

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

func handleMessage(d []byte) error {
	data := map[string]any{}
	err := json.Unmarshal(d, &data)

	if err != nil {
		return err
	}

	return err
}
