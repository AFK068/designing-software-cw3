package consumer

import (
	"context"
	"encoding/json"

	"github.com/AFK068/designing-software-cw3/payment/internal/application/services"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderEvent struct {
	Event       string `json:"event"`
	OrderID     string `json:"order_id"`
	UserID      string `json:"user_id"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

func HandleOrderMessage(paymentService *services.PaymentService, logger *zap.Logger) func(msg *sarama.ConsumerMessage) error {
	return func(msg *sarama.ConsumerMessage) error {
		logger.Info("HandleOrderMessage called", zap.ByteString("raw", msg.Value))

		var event OrderEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			logger.Error("failed to unmarshal order event", zap.Error(err))
			return err
		}

		userID, err := uuid.Parse(event.UserID)
		if err != nil {
			logger.Error("invalid user_id in order event", zap.String("user_id", event.UserID), zap.Error(err))
			return err
		}

		ctx := context.Background()

		_, err = paymentService.ReplenishAccount(ctx, userID, -event.Amount)
		if err != nil {
			logger.Error("failed to debit amount from user", zap.String("user_id", event.UserID), zap.Error(err))
			return err
		}

		logger.Info("amount debited from user", zap.String("user_id", event.UserID), zap.Int64("amount", event.Amount))

		return nil
	}
}
