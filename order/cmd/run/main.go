package main

import (
	"context"

	"github.com/AFK068/designing-software-cw3/order/internal/application/services"
	"github.com/AFK068/designing-software-cw3/order/internal/config"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/httpapi/orderapi"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/kafka/producer"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/outbox"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/server"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/store"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/store/txs"
	"github.com/AFK068/designing-software-cw3/order/internal/migration"
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
			services.NewOrderService,

			// Kafka producer.
			func(cfg *config.Config) (*producer.KafkaAsyncProducer, error) {
				return producer.NewKafkaAsyncProducer([]string{cfg.Kafka.BrokerAddress})
			},

			// Handler.
			orderapi.NewOrderHandler,

			// Outbox worker.
			outbox.NewWorker,

			// Order server.
			server.NewOrderServer,
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
			// Run scrapper.
			func(s *server.OrderServer, lc fx.Lifecycle, logger *zap.Logger) {
				s.RegisterHooks(lc, logger)
			},
		),
	).Run()
}
