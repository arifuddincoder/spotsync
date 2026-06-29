package reservation

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

	authMW := middlewares.AuthMiddleware(jwtService)

	api := e.Group("/api/v1/reservations")

	api.POST("", h.Reserve, authMW)
	api.GET("/my-reservations", h.GetMyReservations, authMW)
	api.DELETE("/:id", h.CancelReservation, authMW)

	api.GET("", h.GetAllReservations, authMW, middlewares.RequireRole("admin"))
}
