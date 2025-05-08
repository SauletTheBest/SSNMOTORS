package http

import (
	"AP2Assignment2/inventory-service/internal/adapter/http/handler"
)

type ProductHandler interface {
	handler.ProductUseCase
}
