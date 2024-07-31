package main

import (
	"api-geteway/api"
	"api-geteway/api/handler"
	"api-geteway/config"
	"api-geteway/logs"
	"api-geteway/queue/kafka/producer"
	"api-geteway/service"

	"github.com/gin-gonic/gin"
)

func main() {
	logger := logs.InitLogger()
	logger.Info("Application started")

	enforcer, err := config.CasbinEnforcer(logger)
	if err != nil {
		logger.Error("Error initializing enforcer", "error", err.Error())
		return
	}

	config := config.Load()
	serviceManager, err := service.NewServiceManager(config)
	if err != nil {
		logger.Error("Error initializing service manager", "error", err.Error())
		return
	}

	writer := producer.NewKafkaProducer([]string{"localhost:9092"})
	defer writer.Close()
	
	handler := handler.NewMainHandler(serviceManager, logger, writer)
	controller := api.NewController(gin.Default())
	controller.SetupRoutes(handler, logger, enforcer)
	controller.StartServer(config)
}
