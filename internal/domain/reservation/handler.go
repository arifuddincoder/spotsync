package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) Reserve(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user context",
		})
	}

	var req dto.CreateReservationRequest
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

	resp, err := h.service.Reserve(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, ErrZoneNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Failed to create reservation",
				Errors:  err.Error(),
			})
		case errors.Is(err, ErrZoneFull), errors.Is(err, ErrDuplicatePlate):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Success: false,
				Message: "Failed to create reservation",
				Errors:  err.Error(),
			})
		default:
			return httpresponse.InternalError(c, "Failed to create reservation", err)
		}
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "Reservation confirmed successfully",
		Data:    resp,
	})
}

func (h *handler) GetMyReservations(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user context",
		})
	}

	resp, err := h.service.GetMyReservations(userID)
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve reservations", err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "My reservations retrieved successfully",
		Data:    resp,
	})
}

func (h *handler) GetAllReservations(c *echo.Context) error {
	resp, err := h.service.GetAllReservations()
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve reservations", err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Reservations retrieved successfully",
		Data:    resp,
	})
}

func (h *handler) CancelReservation(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Success: false,
			Message: "Unauthorized",
			Errors:  "missing user context",
		})
	}
	role, _ := c.Get("user_role").(string)

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid reservation id",
			Errors:  "reservation id must be a positive integer",
		})
	}

	err = h.service.CancelReservation(userID, role, uint(id))
	if err != nil {
		switch {
		case errors.Is(err, ErrReservationNotFound):
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Reservation not found",
				Errors:  err.Error(),
			})
		case errors.Is(err, ErrNotOwner):
			return c.JSON(http.StatusForbidden, httpresponse.Error{
				Success: false,
				Message: "Forbidden",
				Errors:  err.Error(),
			})
		case errors.Is(err, ErrAlreadyCancelled):
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Success: false,
				Message: "Failed to cancel reservation",
				Errors:  err.Error(),
			})
		default:
			return httpresponse.InternalError(c, "Failed to cancel reservation", err)
		}
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Reservation cancelled successfully",
	})
}
