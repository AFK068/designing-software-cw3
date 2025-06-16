package main

import (
	"context"

	"github.com/AFK068/designing-software-cw3/payment/internal/application/services"
	"github.com/AFK068/designing-software-cw3/payment/internal/config"
	"github.com/AFK068/designing-software-cw3/payment/internal/domain"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/httpapi/paymentapi"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/kafka/consumer"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/server"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/store"
	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/store/txs"
	"github.com/AFK068/designing-software-cw3/payment/internal/migration"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	DevConfigPath  = "config/dev.yaml"
	MigrationsPath = "db/migrations"
)

func postgresDB(cfg *config.Config, log *zap.Logger, lc fx.Lifecycle) (domain.Store, *pgxpool.Pool, error) {
	dbPool, err := pgxpool.New(context.Background(), cfg.GetPostgresConnectionString())
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			dbPool.Close()
			return nil
		},
	})

	return store.NewStore(dbPool), dbPool, nil
}

func main() {
	fx.New(
		fx.Provide(
			// Provide logger.
			func() *zap.Logger {
				logger, err := zap.NewProduction()
				if err != nil {
					panic(err)
				}

				return logger
			},

			// Config.
			func() (*config.Config, error) {
				return config.NewConfig(DevConfigPath)
			},

			// Provide database.
			postgresDB,

			// Provide transactor.
			fx.Annotate(
				txs.NewTxBeginner,
				fx.As(new(services.Transactor)),
			),

			// Provide services.
			services.NewPaymentService,

			// Handler.
			paymentapi.NewPaymentHandler,

			// Order server.
			server.NewPaymentServer,

			// Kafka consumer group.
			func(cfg *config.Config, logger *zap.Logger) (*consumer.Group, error) {
				brokers := []string{cfg.Kafka.BrokerAddress}
				groupID := "payment-service"
				return consumer.NewGroup(brokers, groupID, logger)
			},
		),
		fx.Invoke(
			func(cfg *config.Config, log *zap.Logger) {
				if err := migration.RunMigrations(
					cfg.GetPostgresConnectionString(),
					MigrationsPath,
					log,
				); err != nil {
					log.Fatal("failed to run migrations", zap.Error(err))
				}
			},
			func(lc fx.Lifecycle, group *consumer.Group, paymentService *services.PaymentService, logger *zap.Logger) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							topics := []string{"orders.events"}
							for {
								if ctx.Err() != nil {
									return
								}
								err := group.Consume(ctx, topics, consumer.HandleOrderMessage(paymentService, logger))
								if err != nil {
									logger.Error("consumer exited with error", zap.Error(err))
								}
							}
						}()
						return nil
					},
					OnStop: func(_ context.Context) error {
						if err := group.Group.Close(); err != nil {
							logger.Error("failed to close consumer group", zap.Error(err))
						}
						return nil
					},
				})
			},
		),
	).Run()
}
