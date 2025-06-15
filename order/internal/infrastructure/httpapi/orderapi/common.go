package orderapi

import (
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/labstack/echo/v4"

	ordertypes "github.com/AFK068/designing-software-cw3/order/internal/openapi/order/v1"
)

func SendBadRequestResponse(ctx echo.Context, description string) error {
	return ctx.JSON(http.StatusBadRequest, ordertypes.ApiErrorResponse{
		Message: aws.String(description),
		Code:    aws.String("400"),
	})
}

func SendNotFoundResponse(ctx echo.Context, description string) error {
	return ctx.JSON(http.StatusNotFound, ordertypes.ApiErrorResponse{
		Message: aws.String(description),
		Code:    aws.String("404"),
	})
}

func SendInternalErrorResponse(ctx echo.Context, description string) error {
	return ctx.JSON(http.StatusInternalServerError, ordertypes.ApiErrorResponse{
		Message: aws.String(description),
		Code:    aws.String("500"),
	})
}
