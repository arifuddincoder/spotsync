package httpresponse

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v5"
)

type Success struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Error struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func InternalError(c *echo.Context, message string, err error) error {
	log.Printf("[ERROR] %s: %v", message, err)
	return c.JSON(http.StatusInternalServerError, Error{
		Success: false,
		Message: message,
		Errors:  "internal server error",
	})
}
