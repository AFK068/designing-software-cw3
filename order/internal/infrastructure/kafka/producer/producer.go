package producer

import (
	"fmt"

	"github.com/IBM/sarama"
)

type KafkaAsyncProducer struct {
	producer sarama.AsyncProducer
	errors   chan error
}

func NewKafkaAsyncProducer(brokers []string) (*KafkaAsyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Producer.Return.Successes = false
	cfg.Producer.Return.Errors = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 10

	p, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return nil, err
	}

	ap := &KafkaAsyncProducer{
		producer: p,
		errors:   make(chan error, 1),
	}

	go func() {
		for err := range p.Errors() {
			ap.errors <- err.Err
		}
	}()

	return ap, nil
}

func (k *KafkaAsyncProducer) Send(messages ...*sarama.ProducerMessage) error {
	for _, msg := range messages {
		select {
		case k.producer.Input() <- msg:
		case err := <-k.errors:
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}

func (k *KafkaAsyncProducer) Close() error {
	return k.producer.Close()
}
