package phandler

import (
	"api-geteway/api/token"
	"api-geteway/generated/reviews"
	"api-geteway/models"
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --------------------------------------- REVIEWS ------------------------------------------------------

// @Summary Get all reviews
// @Description Get all reviews for a product
// @Tags Review
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Success 200 {object} reviews.GetAllReviewsRequest
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews [get]
func (h *productHandlerIml) GetAllReviews(ctx *gin.Context) {
	var req reviews.GetAllReviewsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		h.logger.Error("ShouldBindQuery", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.GetAllReviews(context.Background(), &req)
	if err != nil {
		h.logger.Error("GetAllReviews", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summar Review by ProductId
// @Description Get reviews by ProductId
// @Tags Review
// @Accept json
// @Produce json
// @Param product_id path string true "Product ID"
// @Param Page query int true "Page"
// @Param Limit query int true "Limit"
// @Success 200 {object} reviews.GetReviewsByPIdResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/{product_id} [get]
func (h *productHandlerIml) GetReviewByProductId(ctx *gin.Context) {
	id := ctx.Param("product_id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	offset, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit"), 10, 64)
	res, err := h.productClient.GetReviewsByProductId(context.Background(), &reviews.GetReviewsByPIdRequest{
		ProductId: id,
		Offset:    offset,
		Limit:     limit,
	})
	if err != nil {
		h.logger.Error("GetReviewsByProductId", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create Review
// @Description Create a new review
// @tags Review
// @Accept json
// @Produce json
// @Param review body models.CreateReview true "Review to create"
// @Param product_id path string true "Product ID"
// @Success 201 {object} reviews.CreateReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/{product_id} [post]
func (h *productHandlerIml) PostReviewByProductId(ctx *gin.Context) {
	id := ctx.Param("product_id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}

	var req models.CreateReview
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
	}
	claims, ok := value.(*token.Claims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
		return
	}

	res, err := h.productClient.CreateReview(context.Background(), &reviews.CreateReviewRequest{
		ProductId: id,
		UserId:    claims.ID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		h.logger.Error("CreateReview", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Create Review by ProductId
// @Description Create a new review for a product
// @tags Review
// @Accept json
// @Produce json
// @param review body reviews.CreateReviewRequest true "Review to create"
// @Success 201 {object} reviews.CreateReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews [post]
func (h *productHandlerIml) PostReview(ctx *gin.Context) {
	var req reviews.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.productClient.CreateReview(ctx, &req)
	if err != nil {
		h.logger.Error("Error in create review", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// @Summary Update Review
// @Description Update an existing review
// @tags Review
// @Accept json
// @Produce json
// @param id path string true "Review ID"
// @param review body models.UpdateReview true "Review to update"
// @Success 200 {object} reviews.UpdateReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/{id} [put]
func (h *productHandlerIml) PutReview(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}

	value, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
	}
	claims, ok := value.(*token.Claims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
		return
	}

	var req models.UpdateReview
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.UpdateReview(context.Background(), &reviews.UpdateReviewRequest{
		Id:        id,
		UserId:    claims.ID,
		ProductId: req.ProductId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		h.logger.Error("UpdateReview", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Update Review
// @Description Update an existing review
// @tags Review
// @Accept json
// @Produce json
// @param id path string true "Review ID"
// @param review body models.UpdateReviewById true "Review to update"
// @Success 200 {object} reviews.UpdateReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/admin/{id} [put]
func (h *productHandlerIml) PutReviewById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}

	var req models.UpdateReviewById
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Error("ShouldBindJSON", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.productClient.UpdateReview(context.Background(), &reviews.UpdateReviewRequest{
		Id:        id,
		UserId:    req.UserId,
		ProductId: req.ProductId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	})
	if err != nil {
		h.logger.Error("UpdateReview", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Review
// @Description Delete an existing review
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} reviews.DeleteReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/{id} [delete]
func (h *productHandlerIml) DeleteReview(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	value, ok := ctx.Get("claims")
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
	}
	claims, ok := value.(*token.Claims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting claims",
		})
		return
	}

	res, err := h.productClient.DeleteReview(context.Background(), &reviews.DeleteReviewRequest{
		Id:     id,
		UserId: claims.ID,
	})
	if err != nil {
		h.logger.Error("DeleteReview", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

// @Summary Delete Review
// @Description Delete an existing review
// @Tags Review
// @Accept json
// @Produce json
// @Param id path string true "Review ID"
// @Success 200 {object} reviews.DeleteReviewResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/reviews/admin/{id} [delete]
func (h *productHandlerIml) DeleteReviewById(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.logger.Error("ProductId not found")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ProductId not found"})
		return
	}
	req := reviews.DeleteReviewRequest{Id: id}

	res, err := h.productClient.DeleteReview(context.Background(), &req)
	if err != nil {
		h.logger.Error("DeleteReview", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}
