package user

import (
	"spotsync/internal/auth"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	repo := NewRepository(db)
	svc := NewService(repo, jwtService)
	h := NewHandler(svc)

	api := e.Group("/api/v1/auth")
	api.POST("/register", h.RegisterUser)
	api.POST("/login", h.LoginUser)
}
