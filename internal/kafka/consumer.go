package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer[T any] struct {
	reader      *kafka.Reader
	logger      *slog.Logger
	workerCount int
}

func NewConsumer[T any](brokers []string, topic, groupID string, workerCount int, logger *slog.Logger) *Consumer[T] {
	if workerCount <= 0 {
		workerCount = 1
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	return &Consumer[T]{
		reader:      r,
		logger:      logger.With("component", "kafka_consumer", "topic", topic, "group", groupID),
		workerCount: workerCount,
	}
}

type message[T any] struct {
	key string
	val T
	msg kafka.Message
}

func (c *Consumer[T]) Close() error {
	return c.reader.Close()
}

func (c *Consumer[T]) Consume(
	ctx context.Context,
	handler func(ctx context.Context, key string, value T, msg kafka.Message) error,
) error {
	ch := make(chan message[T], c.workerCount*2)
	var wg sync.WaitGroup

	for i := 0; i < c.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for m := range ch {
				if err := handler(ctx, m.key, m.val, m.msg); err != nil {
					c.logger.Error("handler error",
						"worker", workerID,
						"key", m.key,
						"error", err,
					)
				}
			}
		}(i)
	}

readLoop:
	for {
		if ctx.Err() != nil {
			break readLoop
		}

		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				break readLoop
			}
			c.logger.Error("failed to read kafka message", "error", err)
			continue
		}

		var value T
		if err := json.Unmarshal(msg.Value, &value); err != nil {
			c.logger.Error("failed to unmarshal kafka message", "error", err)
			continue
		}

		m := message[T]{
			key: string(msg.Key),
			val: value,
			msg: msg,
		}

		select {
		case <-ctx.Done():
			break readLoop
		case ch <- m:
		}
	}

	// Graceful shutdown
	close(ch)
	wg.Wait()
	return ctx.Err()
}
