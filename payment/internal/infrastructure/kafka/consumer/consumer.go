package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type HandleMessageFn func(msg *sarama.ConsumerMessage) error

type Group struct {
	Group   sarama.ConsumerGroup
	groupID string
	logger  *zap.Logger
}

func NewGroup(addrs []string, groupID string, logger *zap.Logger) (*Group, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Group{
		Group:   consumer,
		groupID: groupID,
		logger:  logger,
	}, nil
}

func (c *Group) Consume(ctx context.Context, topics []string, handleMessageFn HandleMessageFn) error {
	consumerHandler := &defaultConsumerHandler{
		ctx:             ctx,
		handleMessageFn: handleMessageFn,
		zapLogger:       c.logger,
	}

	err := c.Group.Consume(ctx, topics, consumerHandler)
	if err != nil {
		c.logger.Error("failed to consume messages", zap.String("group_id", c.groupID), zap.Strings("topics", topics), zap.Error(err))
		return err
	}

	return nil
}

func (c *Group) OnConsume(m *sarama.ConsumerMessage) {
	c.logger.Info(
		"message consumed",
		zap.String("topic", m.Topic),
		zap.Int32("partition", m.Partition),
		zap.Int64("offset", m.Offset),
		zap.ByteString("key", m.Key),
		zap.ByteString("value", m.Value),
	)
}
