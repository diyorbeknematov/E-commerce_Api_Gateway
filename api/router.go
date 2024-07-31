package api

import (
	"api-geteway/api/handler"
	"api-geteway/api/middleware"
	"api-geteway/config"
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	_ "api-geteway/api/handler/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Controller interface {
	SetupRoutes(handler.MainHandler, *slog.Logger, *casbin.Enforcer)
	StartServer(config.Config) error
}

type controllerImpl struct {
	Port   string
	router *gin.Engine
}

func NewController(router *gin.Engine) Controller {
	return &controllerImpl{router: router}
}

func (c *controllerImpl) StartServer(cfg config.Config) error {
	c.Port = cfg.HTTP_PORT
	return c.router.Run(c.Port)
}

// @title Api Getaway
// @version 1.0
// @description api gateway service
// @host localhost:8080
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
// @BasePath /
// @schemes http
func (c *controllerImpl) SetupRoutes(handler handler.MainHandler, logger *slog.Logger, enforcer *casbin.Enforcer) {
	// Implement routes setup here
	c.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router := c.router.Group("/api")
	router.Use(middleware.IsAuthenticated(), middleware.LogMiddleware(logger), middleware.Authorize(enforcer))

	huser := handler.NewUserHandler()
	hproduct := handler.NewProductHandler()

	// ------------------------------- user service --------------------------------
	users := router.Group("/users")
	{
		{ // ---------  user ----------
			users.GET("", huser.GetUser) //ok
			users.PUT("", huser.UpdateUser) //ok
			users.DELETE("", huser.DeleteUser) //ok
			users.GET("/recommendation", huser.GetUserRecommendations) //ok
			users.GET("/products", huser.GetUserOrders)               // ok
		}
		{ // ----------  admin ----------
			users.GET("/:id", huser.GetUserByID) //ok
			users.PUT("/:id", huser.UpdateUserByID) //ok
			users.DELETE("/:id", huser.DeleteUserById) //ok
			users.POST("", huser.CreateUser) //ok
			users.GET("/products/:id", huser.GetUserBoughtProducts) //ok
			users.GET("/list", huser.GetUsers) //ok
		}
	}

	// ----------------------------------- product service -----------------------------------------
	{ // -----------  user ----------
		router.POST("/media", hproduct.PostMedia) //ok
		router.POST("/orders/:product_id", hproduct.PostOrder) // tekshir ok
	}

	basket := router.Group("/basket")
	{
		// -----------  user ----------
		basket.POST("/:product_id", hproduct.AddToBasket) //ok
		basket.GET("", hproduct.GetProductFromBasket) //ok
		basket.DELETE("/:product_id", hproduct.DeleteProductFromBasket) //ok
	}

	products := router.Group("/products")
	{
		{ // -----------  admin ----------
			products.GET("/:id", hproduct.GetProductsByID) //ok
			products.POST("", hproduct.CreateProducts) //ok
			products.PUT("/:id", hproduct.PutProductsByID) //ok
			products.DELETE("/:id", hproduct.DeleteProductsByID) //ok
			products.GET("/order/product_id", hproduct.GetOrdersByProductId)
		}
		{ // -----------  admin and user ----------
			products.GET("/list", hproduct.GetAllProducts) //ok
		}
	}

	categories := router.Group("/categories")
	{
		{ // -----------  admin ----------
			categories.POST("", hproduct.CreateCategories) //ok
			categories.PUT("/:id", hproduct.PutCategoryByID) //ok
			categories.DELETE("/:id", hproduct.DeleteCategoryByID) //ok
		}
		{ // -----------  admin and user ----------
			categories.GET("", hproduct.GetCategories) //ok
		}
	}

	reviews := router.Group("/reviews")
	{
		{ // -----------  user ----------
			reviews.GET("/:product_id", hproduct.GetReviewByProductId) //ok
			reviews.POST("/:product_id", hproduct.PostReviewByProductId)//ok
			reviews.PUT("/:id", hproduct.PutReview) //ok
			reviews.DELETE("/:id", hproduct.DeleteReview) //ok
		}
		{ // -----------  admin ----------
			reviews.POST("", hproduct.PostReview)
			reviews.GET("", hproduct.GetAllReviews)
			reviews.PUT("/admin/:id", hproduct.PutReviewById)
			reviews.DELETE("/admin/:id", hproduct.DeleteReviewById)
		}
	}
}
