package zone

import (
	"spotsync/internal/auth"
	"spotsync/internal/middlewares"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	api := e.Group("/api/v1/zones")

	api.GET("", h.GetAllZones)
	api.GET("/:id", h.GetZone)

	api.POST("", h.CreateZone,
		middlewares.AuthMiddleware(jwtService),
		middlewares.RequireRole("admin"),
	)
}
