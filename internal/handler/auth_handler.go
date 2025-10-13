package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/service"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	response, err := h.authService.Login(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	user, err := h.authService.Register(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
	}

	var req dto.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	err := h.authService.ChangePassword(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to change password",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Password changed successfully",
	})
}

// UserHandler handles user related requests
type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// Create handles user creation
func (h *UserHandler) Create(c echo.Context) error {
	var req dto.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	user, err := h.userService.Create(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetByID handles getting user by ID
func (h *UserHandler) GetByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// Update handles user update
func (h *UserHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	user, err := h.userService.Update(uint(id), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// Delete handles user deletion
func (h *UserHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID",
			Error:   err.Error(),
		})
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// List handles user listing with pagination
func (h *UserHandler) List(c echo.Context) error {
	var params dto.UserQueryParams
	if err := c.Bind(&params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Struct(params); err != nil {
		return c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	users, meta, err := h.userService.List(params)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Failed to retrieve users",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
		Meta:    *meta,
	})
}
