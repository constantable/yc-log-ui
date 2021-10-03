package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/constantable/yc-log-ui/authtoken"
	"github.com/constantable/yc-log-ui/db"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type Handlers struct{}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var dbHandle *db.DB

func (h *Handlers) SetDatabase(db *db.DB) {
	dbHandle = db
}

func (h *Handlers) Login(c echo.Context) error {
	u := new(LoginRequest)
	if err := c.Bind(u); err != nil {
		return echo.ErrBadRequest
	}

	// Check in your db if the user exists or not
	if dbHandle.IsValidUser(u.Username, u.Password) {
		user, err := dbHandle.GetUser(u.Username)
		tokens, err := authtoken.GenerateTokenPair(user, dbHandle)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, tokens)
	}

	return echo.ErrUnauthorized
}

func (h *Handlers) Token(c echo.Context) error {
	type tokenReqBody struct {
		RefreshToken string `json:"refresh_token"`
	}
	tokenReq := tokenReqBody{}

	if err := c.Bind(&tokenReq); err != nil {
		return echo.ErrBadRequest
	}
	if len(tokenReq.RefreshToken) == 0 {
		return echo.ErrBadRequest
	}

	token, err := jwt.Parse(tokenReq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if usr, err := dbHandle.GetUserByRT(int64(claims["sub"].(float64))); err == nil {
			newTokenPair, err := authtoken.GenerateTokenPair(usr, dbHandle)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, newTokenPair)
		}

		return echo.ErrUnauthorized
	}
	return c.JSON(http.StatusForbidden, map[string]string{
		"message": "Expired or invalid refresh token",
	})

}

func (h *Handlers) Password(c echo.Context) error {
	type passwordReqBody struct {
		Password    string `json:"password"`
		OldPassword string `json:"old_password"`
	}
	passwordReq := passwordReqBody{}

	if err := c.Bind(&passwordReq); err != nil {
		return echo.ErrBadRequest
	}
	if len(passwordReq.Password) < 4 {
		return echo.ErrBadRequest
	}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if dbHandle.IsValidUser(claims["username"].(string), passwordReq.OldPassword) {
		usr, err := dbHandle.SetPassword(claims["username"].(string), passwordReq.Password)
		if err != nil {
			log.Println(err)
			return err
		}
		newTokenPair, err := authtoken.GenerateTokenPair(usr, dbHandle)
		if err != nil {
			log.Println(err)
			return err
		}

		return c.JSON(http.StatusOK, newTokenPair)
	}
	return echo.ErrForbidden
}
