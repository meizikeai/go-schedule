package models

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go-schedule/libs/tool"

	"github.com/IBM/sarama"
)

type ModelsKafka struct{}

func NewModelsKafka() *ModelsKafka {
	return &ModelsKafka{}
}

func (m *ModelsKafka) SendKafkaProducerMessage(broker, topic, data string) {
	producer := tools.GetKafkaProducerClient(broker)

	if producer == nil {
		fmt.Println("kafka producer is not connected")
		return
	}

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   nil,
		Value: sarama.StringEncoder(data),
	}

	producer.Input() <- message
}

func (m *ModelsKafka) HandlerKafkaConsumerMessage(broker, topic string) {
	consumer := tools.GetKafkaConsumerClient(broker)
	partitionList, err := consumer.Partitions(topic)

	if err != nil {
		fmt.Printf("error get partition: %v\n", err)
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
func (m *ModelsKafka) ReadKafkaConsumerGroupMessage(callback func([]byte) error) {
	keepRunning := true
	tools.Stdout("Starting a new Sarama consumer")

	// Setup a new Sarama consumer group
	consumer := tool.NewKafkaConsumerGroupConsumer(callback)

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
			if err := client.Consume(ctx, []string{"nyx_traffic"}, consumer); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}

			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
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
