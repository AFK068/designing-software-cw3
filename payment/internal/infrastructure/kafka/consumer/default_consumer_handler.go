package consumer

import (
	"context"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type defaultConsumerHandler struct {
	ctx             context.Context
	handleMessageFn HandleMessageFn
	zapLogger       *zap.Logger
}

func (h *defaultConsumerHandler) Setup(sess sarama.ConsumerGroupSession) error {
	h.zapLogger.Info("consumer group session [setup]", zap.Any("claims", sess.Claims()))
	return nil
}

func (h *defaultConsumerHandler) Cleanup(sess sarama.ConsumerGroupSession) error {
	h.zapLogger.Info("consumer group session [cleanup]", zap.Any("claims", sess.Claims()))
	return nil
}

func (h *defaultConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.zapLogger.Info(
		"ConsumeClaim started",
		zap.Int32("partition", claim.Partition()),
		zap.Int64("initial_offset", claim.InitialOffset()),
		zap.Int64("high_water_mark_offset", claim.HighWaterMarkOffset()),
	)

	defer func() {
		if e := recover(); e != nil {
			h.zapLogger.Error("panic occurred while consuming messages", zap.Any("error", e))
			_ = h.ConsumeClaim(sess, claim)
		}
	}()

	for {
		select {
		case <-h.ctx.Done():
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			h.zapLogger.Info(
				"Message received in ConsumeClaim",
				zap.String("topic", msg.Topic),
				zap.Int32("partition", msg.Partition),
				zap.Int64("offset", msg.Offset),
				zap.ByteString("key", msg.Key),
				zap.ByteString("value", msg.Value),
			)

			err := h.handleMessageFn(msg)
			if err != nil {
				h.zapLogger.Error("failed to handle message", zap.Error(err))
				continue
			}

			sess.MarkMessage(msg, "")
		}
	}
}
