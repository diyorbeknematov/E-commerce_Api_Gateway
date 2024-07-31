package phandler

import (
	"api-geteway/generated/categories"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ---------------------------------- CATIGORIES ----------------------------------------------

// @Summary Get All Categories
// @Description Get all categories
// @Tags Category
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Limit number of categories per page"
// @Success 200 {object} categories.GetAllCategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/categories [get]
func (h *productHandlerIml) GetCategories(ctx *gin.Context) {
	var req categories.GetAllCategoryRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("ShouldBindQuery", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.GetAllCategories(ctx, &req)
	if err != nil {
		h.logger.Error("GetAllCategories", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create Category
// @Description Create a new category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body categories.CreateCategoryRequest true "Category to create"
// @Success 201 {object} categories.CreateCategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/categories [post]
func (h *productHandlerIml) CreateCategories(ctx *gin.Context) {
	var req categories.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.CreateCategory(context.Background(), &req)
	if err != nil {
		h.logger.Error("CreateCategory", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Category
// @Description Update an existing category
// @Tags Category
// @Accept json
// @Produce json
// @Param category body categories.UpdateCategoryRequest true "Category to update"
// @Param id path string true "Category ID"
// @Success 200 {object} categories.UpdateCategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/categories/{id} [put]
func (h *productHandlerIml) PutCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}

	var req categories.UpdateCategoryRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Id = id

	res, err := h.productClient.UpdateCategory(context.Background(), &req)
	if err != nil {
		h.logger.Error("UpdateCategory", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Category
// @Description Delete an existing category
// @Tags Category
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} categories.DeleteCategoryResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/categories/{id} [delete]
func (h *productHandlerIml) DeleteCategoryByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	req := categories.DeleteCategoryRequest{Id: id}
	res, err := h.productClient.DeleteCategory(context.Background(), &req)
	if err != nil {
		h.logger.Error("DeleteCategory", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
