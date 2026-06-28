package server

import (
	"fmt"
	"log"
	"net/http"
	"spotsync/internal/auth"
	"spotsync/internal/config"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}
	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	jwtService, err := auth.NewJWTService(cfg.JwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Hello from go SpotSync")
	})

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
