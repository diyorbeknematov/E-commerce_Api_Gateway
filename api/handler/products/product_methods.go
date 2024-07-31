package phandler

import (
	pbp "api-geteway/generated/products"
	"api-geteway/models"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------- PRODUCTS ----------------------------------------------

// @Summary Create Product
// @Description Create a new product
// @tags Product
// @Accept json
// @Produce json
// @Param product body products.CreateProductRequest true "Create product"
// @Success 201 {object} products.CreateProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products [post]
func (h *productHandlerIml) CreateProducts(ctx *gin.Context) {
	var req = pbp.CreateProductRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.CreateProduct(context.Background(), &req)
	if err != nil {
		h.logger.Error("CreateProducts", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Product
// @Description Update an existing product
// @tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.UpdateProduct true "Update product"
// @Success 200 {object} products.UpdateProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products/{id} [put]
func (h *productHandlerIml) PutProductsByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}

	var req models.UpdateProduct
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	res, err := h.productClient.UpdateProduct(context.Background(), &pbp.UpdateProductRequest{
		Id: id,
		Name: req.Name,
		Description: req.Description,
        Images:      req.Images,
        Price:       req.Price,
        Stock:       req.Stock,
        Discount:    &pbp.Discount{
			DiscountPrice: req.Discount.DiscountPrice,
            Status:       req.Discount.Status,
		},
	})
	if err != nil {
		h.logger.Error("UpdateProducts", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get Product by ID
// @Description Get a product by ID
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} products.GetByIdProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products/{id} [get]
func (h *productHandlerIml) GetProductsByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	req := pbp.GetByIdProductRequest{Id: id}

	res, err := h.productClient.GetByIdProduct(context.Background(), &req)
	if err != nil {
		h.logger.Error("GetProducts", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Get All Products
// @Description Get all products
// @Tags Product
// @Accept json
// @Produce json
// @Param Page query int true "Page number"
// @Param Limit query int true "Limit number of products per page"
// @Param Name query string false "Product name"
// @Param Category query string false "Product category"
// @Param Discount query int false "Product discount"
// @Param PriceOrder query float64 false "Product price"
// @Param RatingOrder query int false "Product rating"
// @Param CommentOrder query int false "Product comment"
// @Param Newest query bool false "Product is new"
// @Success 200 {object} products.GetAllProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products/list [get]
func (h *productHandlerIml) GetAllProducts(ctx *gin.Context) {
	var products pbp.GetAllProductRequest

	if err := ctx.ShouldBindQuery(&products); err != nil {
		h.logger.Error("ShouldBindQuery", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.GetAllProduct(context.Background(), &products)
	if err != nil {
		h.logger.Error("GetAllProduct", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Product
// @Description Delete an existing product
// @Tags Product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} products.DeleteProductResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/products/{id} [delete]
func (h *productHandlerIml) DeleteProductsByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	req := pbp.DeleteProductRequest{Id: id}

	res, err := h.productClient.DeleteProduct(context.Background(), &req)
	if err != nil {
		h.logger.Error("DeleteProducts", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
