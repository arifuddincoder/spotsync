package zone

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/zone/dto"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateZone(c *echo.Context) error {
	var req dto.CreateZoneRequest
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

	response, err := h.service.CreateZone(req)
	if err != nil {
		return httpresponse.InternalError(c, "Failed to create parking zone", err)
	}

	return c.JSON(http.StatusCreated, httpresponse.Success{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    response,
	})
}

func (h *handler) GetAllZones(c *echo.Context) error {
	response, err := h.service.GetAllZones()
	if err != nil {
		return httpresponse.InternalError(c, "Failed to retrieve parking zones", err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    response,
	})
}

func (h *handler) GetZone(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone id",
			Errors:  "zone id must be a positive integer",
		})
	}

	response, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return httpresponse.InternalError(c, "Failed to retrieve parking zone", err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    response,
	})
}
func (h *handler) UpdateZone(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone id",
			Errors:  "zone id must be a positive integer",
		})
	}

	var req dto.UpdateZoneRequest
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

	response, err := h.service.UpdateZone(uint(id), req)
	if err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return httpresponse.InternalError(c, "Failed to update parking zone", err)
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone updated successfully",
		Data:    response,
	})
}

func (h *handler) DeleteZone(c *echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Success: false,
			Message: "Invalid zone id",
			Errors:  "zone id must be a positive integer",
		})
	}

	if err := h.service.DeleteZone(uint(id)); err != nil {
		if errors.Is(err, ErrZoneNotFound) {
			return c.JSON(http.StatusNotFound, httpresponse.Error{
				Success: false,
				Message: "Parking zone not found",
				Errors:  err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Success: false,
			Message: "Failed to delete parking zone",
			Errors:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, httpresponse.Success{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}
