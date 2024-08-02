package phandler

import (
	"api-geteway/api/token"
	"api-geteway/generated/products"
	"api-geteway/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary UploadMediaProduct
// @Security ApiKeyAuth
// @Description Api for upload a new photo
// @Tags Product
// @Accept multipart/form-data
// @Param file formData file true "createUserModel"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/media/ [post]
func (h *productHandlerIml) PostMedia(ctx *gin.Context) {
	h.logger.Info("UploadMediaProduct started")
	header, _ := ctx.FormFile("file")

	url := filepath.Join("media/products", header.Filename)

	err := ctx.SaveUploadedFile(header, url)
	if err != nil {
		return
	}

	h.logger.Info("UploadMediaProduct finished successfully")
	ctx.JSON(200, gin.H{
		"message": url,
	})
}

// @Summary Post Order
// @Description Create a new order
// @Tags Order
// @Accept json
// @Produce json
// @Param product_id path string true "productId" "Product ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/orders/{product_id} [post]
func (h *productHandlerIml) PostOrder(ctx *gin.Context) {
	value, ok := ctx.Get("claims")
	if !ok {
		h.logger.Error("Claims err")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims err"})
		return
	}
	claims, ok := value.(*token.Claims)
	if !ok {
		h.logger.Error("Claims err")
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims err"})
        return
	}
	productId := ctx.Param("product_id")
	
	data, err := json.Marshal(&products.OrderRequest{
		ProductId: productId,
		UserId: claims.ID,
	})
	if err != nil {
		h.logger.Error("JSON Marshal error", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.writer.ProducerMessage("order-created", data)
	if err != nil {
		h.logger.Error("ProducerMessage", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
	})
}

// @Summary Get Orders by product id
// @Description Get orders  by product id
// @Tags Order
// @Accept json
// @Produce json
// @Param product_id path string true "productId" "Product ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/orders/{product_id} [get]
func (h *productHandlerIml) GetOrdersByProductId(ctx *gin.Context) {
	id := ctx.Param("product_id")
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	offset, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	resp, err := h.productClient.GetOrderByPId(ctx, &products.GetOrderByPIdRequest{
		ProductId: id,
		Limit: limit,
		Page: offset,
	})
	if err != nil {
		h.logger.Error("Error in get orders by product id", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, resp)
}

// ---------------------------------- BASKETS -------------------------------------------

// @Summary Add to Basket
// @Description Add a product to the basket
// @Tags Basket
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param product body models.AddToBasket true "Add to basket request"
// @Success 200 {object} products.AddToBasketResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/basket/{product_id} [post]
func (h *productHandlerIml) AddToBasket(ctx *gin.Context) {
	value, ok := ctx.Get("claims")
	if !ok {
		h.logger.Error("Claims not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	claims := value.(*token.Claims)

	id := ctx.Param("product_id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	var req models.AddToBasket 
	if err := ctx.ShouldBindJSON(&req); err!= nil {
		h.logger.Error("Bind JSON error", "error", err.Error())
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
	}

	res, err := h.productClient.AddToBasket(context.Background(), &products.AddToBasketRequest{
		ProductId: id,
		UserId: claims.ID,
		PurchaseDate: req.PurchaseDate,
		Quantity: req.Quantity,
        Price: req.Price,
	})
	if err != nil {
		h.logger.Error("AddToBasket", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Basket
// @Description Get products in the basket
// @Tags Basket
// @Accept json
// @Produce json
// @Success 200 {object} products.BasketResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/basket [get]
func (h *productHandlerIml) GetProductFromBasket(ctx *gin.Context) {
	value, ok := ctx.Get("claims")
	if !ok {
		h.logger.Error("Claims not found")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	claims := value.(*token.Claims)

	req := products.GetBasketRequest{UserId: claims.ID}
	products, err := h.productClient.GetBasketProducts(context.Background(), &req)
	if err != nil {
		h.logger.Error("GetBasketProducts", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// @Summary Delete from Basket
// @Description Delete a product from the basket
// @Tags Basket
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} products.DeleteBasketResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/basket/{product_id} [delete]
func (h *productHandlerIml) DeleteProductFromBasket(ctx *gin.Context) {
	id := ctx.Param("product_id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	fmt.Print(id)
	fmt.Print(id)
	fmt.Print(id)

	req := products.DeleteBasketRequest{ProductId: id}
	res, err := h.productClient.DeleteBasketProduct(context.Background(), &req)
	if err != nil {
		h.logger.Error("DeleteBasket", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
