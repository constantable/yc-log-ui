package authtoken

import (
	"github.com/constantable/yc-log-ui/db"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateTokenPair(user db.User, dbHandle *db.DB) (map[string]string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Username
	claims["admin"] = user.IsAdmin
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	secret := os.Getenv("SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	tokents := int64(user.Id*1e9) + time.Now().Unix()
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = tokents
	rtClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	rt, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	err = dbHandle.NewRefreshToken(user.Id, tokents)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}
