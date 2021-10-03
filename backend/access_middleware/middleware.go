package access_middleware

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var secret string

func getSecret() string {
	if secret == "" {
		err := godotenv.Load(".env.local")
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		secret = os.Getenv("SECRET")
	}
	return secret
}

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(getSecret()),
})

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["admin"].(bool)

		if !isAdmin {
			return echo.ErrForbidden
		}

		return next(c)
	}
}
