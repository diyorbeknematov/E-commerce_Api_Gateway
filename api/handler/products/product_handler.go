package phandler

import (
	"api-geteway/generated/mainservice"
	"api-geteway/queue/kafka/producer"
	"api-geteway/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type ProductHandler interface {
	// Products
	CreateProducts(ctx *gin.Context)
	PutProductsByID(ctx *gin.Context)
	GetProductsByID(ctx *gin.Context)
	GetAllProducts(ctx *gin.Context)
	DeleteProductsByID(ctx *gin.Context)
	// Basket
	AddToBasket(ctx *gin.Context)
	GetProductFromBasket(ctx *gin.Context)
	DeleteProductFromBasket(ctx *gin.Context)
	// Media
	PostMedia(ctx *gin.Context) 
	// Post Order
	PostOrder(ctx *gin.Context) // NOT WORK
	GetOrdersByProductId(ctx *gin.Context)
	//Catigories
	CreateCategories(ctx *gin.Context)
	PutCategoryByID(ctx *gin.Context)
	GetCategories(ctx *gin.Context)
	DeleteCategoryByID(ctx *gin.Context)
	// Review
	GetReviewByProductId(ctx *gin.Context)
	PostReview(ctx *gin.Context)
	PostReviewByProductId(ctx *gin.Context)
	PutReview(ctx *gin.Context)
	PutReviewById(ctx *gin.Context)
	DeleteReview(ctx *gin.Context)
	DeleteReviewById(ctx *gin.Context)
	GetAllReviews(ctx *gin.Context)
}

type productHandlerIml struct {
	productClient mainservice.MainServiceClient
	logger        *slog.Logger
	writer        producer.KafkaProducer
}

func NewProductHandler(serviceManger service.ServiceManager, logger *slog.Logger, writer producer.KafkaProducer) ProductHandler {
	return &productHandlerIml{
		productClient: serviceManger.ProductService(),
		logger:        logger,
		writer:        writer,
	}
}
