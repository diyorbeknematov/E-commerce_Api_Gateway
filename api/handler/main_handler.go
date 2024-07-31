package handler

import (
	phandler "api-geteway/api/handler/products"
	uhandler "api-geteway/api/handler/user"
	"api-geteway/queue/kafka/producer"
	"api-geteway/service"
	"log/slog"
)

type Handler struct {
	User    uhandler.UserHandler
	Product phandler.ProductHandler
}

func (h *Handler) NewUserHandler() uhandler.UserHandler {
	return h.User
}

func (h *Handler) NewProductHandler() phandler.ProductHandler {
	return h.Product
}

type MainHandler interface {
	NewUserHandler() uhandler.UserHandler
	NewProductHandler() phandler.ProductHandler
}

func NewMainHandler(serviceManger service.ServiceManager, logger *slog.Logger, writer producer.KafkaProducer) MainHandler {
	return &Handler{
		User:    uhandler.NewUserHandler(serviceManger, logger),
		Product: phandler.NewProductHandler(serviceManger, logger, writer),
	}
}
