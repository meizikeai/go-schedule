// internal/pkg/database/kafka/kafka.go
package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-schedule/internal/config"

	"github.com/segmentio/kafka-go"
)

type Client struct {
	instances map[string]config.KafkaInstance
	writers   map[string]*kafka.Writer
	readers   map[string]*kafka.Reader
	mu        sync.Mutex
}

func NewClient(cfgs *map[string]config.KafkaInstance) *Client {
	c := &Client{
		instances: make(map[string]config.KafkaInstance),
		writers:   make(map[string]*kafka.Writer),
		readers:   make(map[string]*kafka.Reader),
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

func readerKey(instance, topic, group string) string {
	if group == "" {
		return fmt.Sprintf("%s:%s:no-group", instance, topic)
	}
	return fmt.Sprintf("%s:%s:group:%s", instance, topic, group)
}

func (c *Client) getReader(instance, topic, group string) (*kafka.Reader, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := readerKey(instance, topic, group)
	if r, ok := c.readers[key]; ok {
		return r, nil
	}

	cfg, ok := c.instances[instance]
	if !ok {
		return nil, fmt.Errorf("kafka instance %s not found", instance)
	}

	rc := kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   topic,
	}

	if group != "" {
		rc.GroupID = group
	}

	r := kafka.NewReader(rc)
	c.readers[key] = r
	return r, nil
}

func (c *Client) Consume(ctx context.Context, instance string, topic string, group string) (kafka.Message, error) {
	_, ok := c.writers[instance]
	if !ok {
		return kafka.Message{}, fmt.Errorf("kafka instance %s not found", instance)
	}

	r, err := c.getReader(instance, topic, group)
	if err != nil {
		return kafka.Message{}, err
	}

	return r.FetchMessage(ctx)
}

func (c *Client) Commit(ctx context.Context, instance, topic, group string, msg kafka.Message) error {
	if group == "" {
		return nil
	}

	key := readerKey(instance, topic, group)

	c.mu.Lock()
	r := c.readers[key]
	c.mu.Unlock()

	if r == nil {
		return nil
	}
	return r.CommitMessages(ctx, msg)
}

func (c *Client) Close() {
	for _, w := range c.writers {
		_ = w.Close()
	}
	for _, r := range c.readers {
		_ = r.Close()
	}
}

// Usage

// 1、Producer（Key = nil）
// 	err := client.Produce(ctx, "instance", "topic", nil, []byte(`{"id":123,"type":"email"}`))

// Consumer（worker group）
// for {
// 	msg, err := client.Consume(ctx, "instance", "topic", "group")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	handleTask(msg.Value)

// 	_ = client.Commit(ctx, "instance", "topic", "group", msg)
// }

// 2、Producer（Key = user_id）
// err := client.Produce(ctx, "instance", "topic", nil, []byte(`{"id":123,"type":"email"}`))

// Consumer（worker group）
// for {
// 	msg, err := client.Consume(ctx, "instance", "topic", "group")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	handleTask(msg.Value)

// 	_ = client.Commit(ctx, "instance", "topic", "group", msg)
// }

// 3、Consumer
// for {
// 	msg, err := client.Consume(ctx, "instance", "topic", "")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	handleTask(msg.Value)
// }

// go func() {
// 	for {
// 		msg, err := client.Consume(ctx, "instance", "topic", "")
// 		if err != nil {
// 			if ctx.Err() != nil {
// 				return
// 			}
// 			continue
// 		}
// 		handleTask(msg.value)
// 	}
// }()
