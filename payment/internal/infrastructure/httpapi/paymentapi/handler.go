package paymentapi

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime/types"

	"github.com/AFK068/designing-software-cw3/payment/internal/application/services"

	paymenttypes "github.com/AFK068/designing-software-cw3/payment/internal/openapi/payment/v1"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler(paymentService *services.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) PostAccount(ctx echo.Context) error {
	var req paymenttypes.PostAccountJSONRequestBody

	if err := ctx.Bind(&req); err != nil || req.UserId == nil {
		return ctx.JSON(400, paymenttypes.ApiErrorResponse{Message: aws.String("invalid request"), Code: aws.String("bad_request")})
	}

	account, err := h.paymentService.CreateAccount(ctx.Request().Context(), *req.UserId)

	if err != nil {
		return ctx.JSON(409, paymenttypes.ApiErrorResponse{Message: aws.String(err.Error()), Code: aws.String("account_exists")})
	}

	return ctx.JSON(201, map[string]interface{}{"user_id": account.UserID, "balance": account.Amount})
}

func (h *PaymentHandler) GetAccountUserIdBalance( //nolint:revive,stylecheck // according to codgen interface
	ctx echo.Context,
	userID types.UUID,
) error {
	account, err := h.paymentService.GetAccount(ctx.Request().Context(), userID)

	if err != nil {
		return ctx.JSON(404, paymenttypes.ApiErrorResponse{Message: aws.String("account not found"), Code: aws.String("not_found")})
	}

	return ctx.JSON(200, paymenttypes.AccountBalanceResponse{
		UserId:  &account.UserID,
		Balance: aws.String(fmt.Sprintf("%d", account.Amount)),
	})
}

func (h *PaymentHandler) PostAccountUserIdDeposit( //nolint:revive,stylecheck // according to codgen interface
	ctx echo.Context,
	userID types.UUID,
) error {
	var req paymenttypes.PostAccountUserIdDepositJSONRequestBody

	if err := ctx.Bind(&req); err != nil || req.Amount == nil {
		return ctx.JSON(400, paymenttypes.ApiErrorResponse{Message: aws.String("invalid request"), Code: aws.String("bad_request")})
	}

	account, err := h.paymentService.ReplenishAccount(ctx.Request().Context(), userID, *req.Amount)

	if err != nil {
		return ctx.JSON(404, paymenttypes.ApiErrorResponse{Message: aws.String("account not found"), Code: aws.String("not_found")})
	}

	return ctx.JSON(200, map[string]interface{}{"user_id": account.UserID, "balance": account.Amount})
}
