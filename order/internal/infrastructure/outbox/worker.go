package outbox

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/kafka/producer"
	"github.com/IBM/sarama"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

const (
	batchSize          = 500
	DefaultJobDuration = 3 * time.Second
)

type Worker struct {
	scheduler gocron.Scheduler
	store     domain.Store
	producer  *producer.KafkaAsyncProducer
	logger    *zap.Logger
}

func NewWorker(store domain.Store, producer *producer.KafkaAsyncProducer, logger *zap.Logger) (*Worker, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Worker{
		scheduler: scheduler,
		store:     store,
		producer:  producer,
		logger:    logger,
	}, nil
}

func (w *Worker) Run(jobDuration time.Duration) {
	_, err := w.scheduler.NewJob(
		gocron.DurationJob(
			jobDuration,
		),
		gocron.NewTask(
			w.ProcessOutbox,
		),
	)

	if err != nil {
		w.logger.Error("failed to create new job in scheduler", zap.Error(err))
		return
	}

	w.logger.Info("outbox worker started", zap.Duration("jobDuration", jobDuration))
	w.scheduler.Start()
}

func (w *Worker) ProcessOutbox() {
	ctx := context.Background()

	outboxRecords, err := w.store.GetUnprocessedOutbox(ctx, batchSize)
	if err != nil {
		w.logger.Error("failed to get unprocessed outbox records", zap.Error(err))
	}

	type userIDPayload struct {
		UserID string `json:"user_id"`
	}

	for _, rec := range outboxRecords {
		var payload userIDPayload
		if err := json.Unmarshal(rec.Payload, &payload); err != nil {
			w.logger.Error("failed to unmarshal outbox payload", zap.Error(err), zap.Any("recordID", rec.ID))
			continue
		}

		msg := &sarama.ProducerMessage{
			Topic: "orders.events",
			Key:   sarama.StringEncoder(payload.UserID),
			Value: sarama.ByteEncoder(rec.Payload),
		}

		err := w.producer.Send(msg)
		if err != nil {
			w.logger.Error("failed to send message to kafka", zap.Error(err), zap.Any("recordID", rec.ID))
			continue
		}

		if err := w.store.MarkOutboxProcessed(ctx, rec.ID); err != nil {
			w.logger.Error("failed to mark outbox record as processed", zap.Error(err), zap.Any("recordID", rec.ID))
		} else {
			w.logger.Info("outbox record processed", zap.Any("recordID", rec.ID))
		}
	}

	w.logger.Info("outbox processing completed", zap.Int("processedRecords", len(outboxRecords)))
}

func (w *Worker) Stop() error {
	err := w.scheduler.Shutdown()
	if err != nil {
		w.logger.Error("failed to stop scheduler")
		return fmt.Errorf("failed to stop scheduler: %w", err)
	}

	w.logger.Info("Scheduler stopped")

	return nil
}
