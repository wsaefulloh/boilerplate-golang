package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wsaefulloh/boilerplate-golang/internal/domain"
	"github.com/wsaefulloh/boilerplate-golang/pkg/response"
)

type UserHandler struct {
	userService domain.UserService
}

func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID format", err.Error())
		return
	}

	user, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to retrieve user")
		return
	}

	response.OK(c, "User retrieved successfully", user)
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, totalItems, err := h.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		response.InternalServerError(c, "Failed to retrieve user list")
		return
	}

	response.Paginated(c, 200, "Users retrieved successfully", users, page, pageSize, totalItems)
}

func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID format", err.Error())
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	updatedUser, err := h.userService.UpdateUser(c.Request.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		if errors.Is(err, domain.ErrEmailAlreadyExists) {
			response.BadRequest(c, "Update failed", "Email address is already in use")
			return
		}
		response.InternalServerError(c, "Failed to update user")
		return
	}

	response.OK(c, "User updated successfully", updatedUser)
}

func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid user ID format", err.Error())
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c, "User not found")
			return
		}
		response.InternalServerError(c, "Failed to delete user")
		return
	}

	response.OK(c, "User deleted successfully", nil)
}
