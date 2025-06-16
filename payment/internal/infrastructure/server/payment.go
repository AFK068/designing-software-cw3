package server

import (
	"context"
	"time"

	"github.com/AFK068/designing-software-cw3/payment/internal/infrastructure/httpapi/paymentapi"
	paymenttypes "github.com/AFK068/designing-software-cw3/payment/internal/openapi/payment/v1"
	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type PaymentServer struct {
	Handler *paymentapi.PaymentHandler
	Echo    *echo.Echo
}

func NewPaymentServer(handler *paymentapi.PaymentHandler) *PaymentServer {
	return &PaymentServer{
		Handler: handler,
		Echo:    echo.New(),
	}
}

func (s *PaymentServer) Start() error {
	paymenttypes.RegisterHandlers(s.Echo, s.Handler)

	return s.Echo.Start(":" + "8083")
}

func (s *PaymentServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *PaymentServer) RegisterHooks(lc fx.Lifecycle, log *zap.Logger) {
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

			return nil
		},
	})
}
