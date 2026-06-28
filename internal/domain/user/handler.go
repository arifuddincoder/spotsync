package user

import (
	"errors"
	"net/http"

	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) RegisterUser(c *echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.RegisterUser(req)
	if err != nil {
		if errors.Is(err, ErrAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Success: false,
				Message: "Failed to create user",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to create user",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "User registered successfully",
		Data:    response,
	})
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid request payload",
			Errors:  err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Validation failed",
			Errors:  err.Error(),
		})
	}

	response, err := h.service.LoginUser(req)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Success: false,
				Message: "Login failed",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to login user",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}
