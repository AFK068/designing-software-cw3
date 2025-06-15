package server

import (
	"context"
	"time"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/httpapi/orderapi"
	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/outbox"
	ordertypes "github.com/AFK068/designing-software-cw3/order/internal/openapi/order/v1"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type OrderServer struct {
	Handler *orderapi.OrderHandler
	Echo    *echo.Echo
	Worker  *outbox.Worker
}

func NewOrderServer(handler *orderapi.OrderHandler, worker *outbox.Worker) *OrderServer {
	return &OrderServer{
		Handler: handler,
		Echo:    echo.New(),
		Worker:  worker,
	}
}

func (s *OrderServer) Start() error {
	ordertypes.RegisterHandlers(s.Echo, s.Handler)

	s.Worker.Run(outbox.DefaultJobDuration)

	return s.Echo.Start(":" + "8082")
}

func (s *OrderServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *OrderServer) RegisterHooks(lc fx.Lifecycle, log *zap.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Info("Starting order server")

			go func() {
				if err := s.Start(); err != nil {
					log.Error("Failed to start order server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			log.Info("Stopping order server")

			if err := s.Stop(); err != nil {
				log.Error("Failed to stop order server", zap.Error(err))
			}

			if err := s.Worker.Stop(); err != nil {
				log.Error("Failed to stop outbox worker", zap.Error(err))
			}

			return nil
		},
	})
}
