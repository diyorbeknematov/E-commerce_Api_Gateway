package uhandler

import (
	"api-geteway/generated/mainservice"
	"api-geteway/generated/user"
	"api-geteway/service"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	CreateUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	UpdateUserByID(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	DeleteUserById(ctx *gin.Context)
	DeleteUser(ctx *gin.Context)
	GetUserRecommendations(ctx *gin.Context)
	GetUserOrders(ctx *gin.Context)
	GetUserBoughtProducts(ctx *gin.Context)
}

type userHandlerIml struct {
	userClient    user.UserServiceClient
	productClient mainservice.MainServiceClient
	logger        *slog.Logger
}

func NewUserHandler(serviceManger service.ServiceManager, logger *slog.Logger) UserHandler {
	return &userHandlerIml{
		userClient:    serviceManger.UserService(),
		productClient: serviceManger.ProductService(),
		logger:        logger,
	}
}
