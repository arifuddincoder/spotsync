package middlewares

import (
	"net/http"
	"strings"

	"spotsync/internal/auth"
	"spotsync/internal/httpresponse"

	"github.com/labstack/echo/v5"
)

func AuthMiddleware(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "missing authorization header",
				})
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "invalid authorization header format",
				})
			}

			claims, err := jwtService.ValidateToken(parts[1])
			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpresponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
			c.Set("user_name", claims.Name)
			c.Set("user_role", claims.Role)
			return next(c)
		}
	}
}

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok || userRole != role {
				return c.JSON(http.StatusForbidden, httpresponse.Error{
					Success: false,
					Message: "Forbidden",
					Errors:  "insufficient permissions",
				})
			}
			return next(c)
		}
	}
}
