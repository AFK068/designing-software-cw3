package orderapi

import (
	"net/http"

	"github.com/AFK068/designing-software-cw3/order/internal/application/mapper"
	"github.com/AFK068/designing-software-cw3/order/internal/application/services"
	ordertypes "github.com/AFK068/designing-software-cw3/order/internal/openapi/order/v1"
	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) PostOrder(ctx echo.Context) error {
	var req ordertypes.PostOrderJSONRequestBody
	if err := ctx.Bind(&req); err != nil || req.UserId == nil || req.Amount == nil || req.Description == nil {
		return SendBadRequestResponse(ctx, "invalid request")
	}

	order, err := h.orderService.CreateOrder(ctx.Request().Context(), *req.UserId, *req.Amount, *req.Description)

	if err != nil {
		return SendInternalErrorResponse(ctx, err.Error())
	}

	return ctx.JSON(http.StatusCreated, mapper.OrderToOrderResponse(order))
}

func (h *OrderHandler) GetOrderOrderIdStatus( //nolint:revive,stylecheck // according to codgen interface
	ctx echo.Context,
	orderID openapi_types.UUID,
) error {
	order, err := h.orderService.GetOrder(ctx.Request().Context(), orderID)

	if err != nil {
		return SendNotFoundResponse(ctx, "order not found")
	}

	return ctx.JSON(http.StatusOK, ordertypes.OrderStatusResponse{
		OrderId: &order.ID,
		Status:  aws.String(string(order.Status)),
	})
}

func (h *OrderHandler) GetOrders(ctx echo.Context, params ordertypes.GetOrdersParams) error {
	orders, err := h.orderService.ListOrders(ctx.Request().Context(), params.UserId)
	if err != nil {
		return SendInternalErrorResponse(ctx, err.Error())
	}

	resp := make([]ordertypes.OrderResponse, 0, len(orders))

	for _, o := range orders {
		resp = append(resp, mapper.OrderToOrderResponse(o))
	}

	return ctx.JSON(http.StatusOK, resp)
}
