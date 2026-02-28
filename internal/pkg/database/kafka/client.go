// internal/pkg/database/kafka/kafka.go
package kafka

import (
	"context"
	"fmt"
	"time"

	"go-schedule/internal/config"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	instances map[string]config.KafkaInstance
	writers   map[string]*kafka.Writer
}

func NewClient(cfgs *map[string]config.KafkaInstance) *Client {
	c := &Client{
		instances: make(map[string]config.KafkaInstance),
		writers:   make(map[string]*kafka.Writer),
	}

	for name, cfg := range *cfgs {
		c.instances[name] = cfg
		c.writers[name] = &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
		}
	}

	return c
}

func (c *Client) Produce(ctx context.Context, instance string, topic string, key, value []byte) error {
	w, ok := c.writers[instance]
	if !ok {
		return fmt.Errorf("kafka writer %s not found", instance)
	}

	return w.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

func (c *Client) NewKafkaReader(instance, topic, group string) (*kafka.Reader, error) {
	cfg, ok := c.instances[instance]
	if !ok {
		return nil, fmt.Errorf("kafka instance %s not found", instance)
	}

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        group,
		Topic:          topic,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 0,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		MaxWait:        500 * time.Millisecond,
	}), nil
}

func (c *Client) Consume(ctx context.Context, instance, topic, group string, handler func(kafka.Message) error) {
	r, err := c.NewKafkaReader(instance, topic, group)
	if err != nil {
		fmt.Printf("Failed to create reader: %v\n", err)
		return
	}
	defer r.Close()

	for {
		msg, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			time.Sleep(time.Second)
			continue
		}

		if err := handler(msg); err != nil {
			if err == context.Canceled {
				return
			}
			fmt.Printf("Failed to process message: %v\n", err)
			continue
		}

		if err := r.CommitMessages(ctx, msg); err != nil {
			fmt.Printf("Failed to commit messages: %v\n", err)
		}
	}
}

func (c *Client) Close() {
	for _, w := range c.writers {
		_ = w.Close()
	}
}

// Usage

// 1、Producer
// 	err := client.Produce(ctx, "instance", "topic", nil, []byte(`{"id":123,"type":"email"}`))

// 2、Consumer
// for i := 0; i < 3; i++ {
// 	go client.Consume(ctx, "instance", "topic", "group", func(m kafka.Message) error {
// 		fmt.Printf("Worker %d processing: %s\n", i, string(m.Value))
// 		return nil
// 	})
// }
