package uhandler

import (
	"api-geteway/api/token"
	"api-geteway/generated/products"
	"api-geteway/generated/user"
	"api-geteway/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get All Users
// @Description Get all users for admin
// @Tags User
// @Accept json
// @Produce json
// @Param FullName query string false "Full Name filter"
// @Param City query string false "City filter"
// @Param State query string false "State filter"
// @Param Country query string false "Country filter"
// @Param Limit query int true "Limit of records"
// @Param Offset query int true "Offset for pagination"
// @Success 200 {object} user.GetAllUsersResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/list [get]
func (h *userHandlerIml) GetUsers(ctx *gin.Context) {
	// Logic for getting all users admin uchun
	var fUser user.GetAllUsersRequest
	if err := ctx.ShouldBindQuery(&fUser); err != nil {
		h.logger.Error("Invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userClient.GetAllUsers(ctx, &fUser)
	if err != nil {
		h.logger.Error("Error getting users", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting users",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get user
// @Description Get user by token
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} user.GetUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [get]
func (h *userHandlerIml) GetUser(ctx *gin.Context) {
	// Logic for getting a user by id

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
	resp, err := h.userClient.GetUser(ctx, &user.GetUserRequest{
		Id: claims.ID,
	})
	if err != nil {
		h.logger.Error("Error getting user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get User By ID
// @Description Get user by Id for admin
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "user id"
// @Success 200 {object} user.GetUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/{id} [get]
func (h *userHandlerIml) GetUserByID(ctx *gin.Context) {
	// Extract the user ID from the URL parameters
	id := ctx.Param("id")

	if id == "" {
		h.logger.Error("Invalid user ID")
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	resp, err := h.userClient.GetUser(ctx, &user.GetUserRequest{Id: id})
	if err != nil {
		h.logger.Error("Error getting user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Create User
// @Description Create user for admin
// @Tags User
// @Accept json
// @Produce json
// @Param CreateUser body user.CreateUserRequest true "create user"
// @Success 200 {object} user.CreateUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [post]
func (h *userHandlerIml) CreateUser(ctx *gin.Context) {
	// Logic for creating a user admin uchun
	var user user.CreateUserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		h.logger.Error("Invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userClient.CreateUser(ctx, &user)
	if err != nil {
		h.logger.Error("Error creating user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error creating user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, resp)
}

// @Summary Update User By ID
// @Description Update user by id for admin
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUser body models.UpdateUserById true "update user"
// @Success 200 {object} user.UpdateUserByIdResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [put]
func (h *userHandlerIml) UpdateUserByID(ctx *gin.Context) {
	// Logic for updating a user by id admin uchun
	id := ctx.Param("id")
	var updateUser models.UpdateUserById
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		h.logger.Error("Invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}
	resp, err := h.userClient.UpdateUserById(ctx, &user.UpdateUserByIdRequest{
		Id:         id,
		FullName:   updateUser.FullName,
		Username:   updateUser.Username,
		Phone:      updateUser.Phone,
		Email:      updateUser.Email,
		Image:      updateUser.Image,
		Role:       updateUser.Role,
		Address:    updateUser.Address,
		City:       updateUser.City,
		State:      updateUser.State,
		Country:    updateUser.Country,
		PostalCode: updateUser.PostalCode,
	})
	if err != nil {
		h.logger.Error("Error updating user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error updating user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Update User
// @Description Update user by token
// @Tags User
// @Accept json
// @Produce json
// @Param UpdateUser body models.UpdateUser true "update user"
// @Success 200 {object} user.UpdateUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [put]
func (h *userHandlerIml) UpdateUser(ctx *gin.Context) {
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

	// Logic for updating a user
	var updateUser models.UpdateUser
	if err := ctx.ShouldBindJSON(&updateUser); err != nil {
		h.logger.Error("Invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}

	resp, err := h.userClient.UpdateUser(ctx, &user.UpdateUserRequest{
		Id:          claims.Id,
		FullName:    updateUser.FullName,
		Username:    updateUser.Username,
		PhoneNumber: updateUser.PhoneNumber,
		Email:       updateUser.Email,
		Image:       updateUser.Image,
		NewPasswrod: updateUser.NewPasswrod,
		Address:     updateUser.Address,
		City:        updateUser.City,
		State:       updateUser.State,
		Country:     updateUser.Country,
		PostalCode:  updateUser.PostalCode,
		Password:    updateUser.Password,
	})
	if err != nil {
		h.logger.Error("Error updating user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error updating user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)

}

// @Summary Delete User By ID
// @Description Delete a user by ID for admin
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} user.DeleteUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/{id} [delete]
func (h *userHandlerIml) DeleteUserById(ctx *gin.Context) {
	// Logic for deleting a user by id admin uchun
	id := ctx.Param("id")

	resp, err := h.userClient.DeleteUserByID(ctx, &user.DeleteByIdRequest{UserId: id})
	if err != nil {
		h.logger.Error("Error deleting user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error deleting user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Delete User
// @Description Delete user by token
// @Tags User
// @Accept json
// @Produce json
// @Param password body models.Authenticated true "User Password"
// @Success 204 {object} user.DeleteUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [delete]
func (h *userHandlerIml) DeleteUser(ctx *gin.Context) {
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

	var p models.Authenticated
	if err := ctx.BindJSON(&p); err != nil {
		h.logger.Error("Invalid request", "error", err.Error())
		ctx.JSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Invalid request",
			Error:   err.Error(),
		})
		return
	}
	resp, err := h.userClient.DeleteUser(ctx, &user.DeleteUserRequest{Id: claims.ID, Password: p.Password})
	if err != nil {
		h.logger.Error("Error deleting user", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error deleting user",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get User Recommendations
// @Description Get user recommendations
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} products.GetRecommendationsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/recommendation [get]
func (h *userHandlerIml) GetUserRecommendations(ctx *gin.Context) {
	// Logic for getting user recommendations
	fmt.Println("hello")
	resp, err := h.productClient.GetUserRecommendation(ctx, &products.Void{})
	if err != nil {
		h.logger.Error("Error getting user recommendations", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting user recommendations",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get User Orders
// @Description Get user's orders
// @tags User
// @Accept json
// @Produce json
// @Param limit query int true "Page Limit"
// @Param offset query int true "Page Offset"
// @Success 200 {object} products.GetPurchasedPResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/products [get]
func (h *userHandlerIml) GetUserOrders(ctx *gin.Context) {
	// Logic for getting user's orders
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

	limit, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	offset, _ := strconv.ParseInt(ctx.Query("offset"), 10, 64)

	resp, err := h.productClient.GetPurchasedProducts(ctx, &products.GetPurchasedPRequest{
		UserId: claims.ID,
		Limit:  limit,
		Page:   offset,
	})
	if err != nil {
		h.logger.Error("Error getting user's orders", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting user's orders",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

// @Summary Get User Bought Products
// @Descripton Get User Bought Products
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param limit query int true "Page Limit"
// @Param page query int true "Page Offset"
// @Success 200 {object} products.GetPurchasedPResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users/products/{id} [get]
func (h *userHandlerIml) GetUserBoughtProducts(ctx *gin.Context) {
	// Logic for getting user's purchased products
	id := ctx.Param("id")
	limit, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	offset, _ := strconv.ParseInt(ctx.Query("offset"), 10, 64)

	resp, err := h.productClient.GetPurchasedProducts(ctx, &products.GetPurchasedPRequest{
		UserId: id,
		Limit:  limit,
		Page:   offset,
	})
	if err != nil {
		h.logger.Error("Error getting user's purchased products", "error", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Error getting user's purchased products",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resp)
}
