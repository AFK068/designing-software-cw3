package mapper

import (
	"github.com/aws/aws-sdk-go-v2/aws"

	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/AFK068/designing-software-cw3/order/internal/infrastructure/domain"

	ordertypes "github.com/AFK068/designing-software-cw3/order/internal/openapi/order/v1"
)

func OrderToOrderResponse(o *domain.Order) ordertypes.OrderResponse {
	return ordertypes.OrderResponse{
		OrderId:     (*openapi_types.UUID)(&o.ID),
		UserId:      (*openapi_types.UUID)(&o.UserID),
		Description: aws.String(o.Description),
		Status:      aws.String(string(o.Status)),
		Amount:      aws.Int64(o.Amount),
	}
}
