.PHONY: generate_openapi
generate_openapi:
	@if ! command -v 'oapi-codegen' &> /dev/null; then \
		echo "Please install oapi-codegen!"; exit 1; \
	fi;

	@mkdir -p payment/internal/openapi/payment/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		openapi/payment/v1/payment-api.yaml > payment/internal/openapi/payment/v1/payment-api.gen.go

	@mkdir -p order/internal/openapi/order/v1
	@oapi-codegen -package v1 \
		-generate server,types \
		openapi/order/v1/order-api.yaml > order/internal/openapi/order/v1/order-api.gen.go
