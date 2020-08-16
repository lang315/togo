package transport

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/lang315/togo/internal/storages"
	"time"
)

const jwtSecret = "AqsA46R925YquUaLvu5mGJNj"

type UserClaims struct {
	UserID string
	jwt.StandardClaims
}

func makeUserJWTToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	})
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return t, nil
}

func getUserID(ctx echo.Context) string {
	u := ctx.Get("user").(*jwt.Token)
	claims := u.Claims.(*UserClaims)
	return claims.UserID
}

func getCurrentUser(db *storages.Data, ctx echo.Context) (*storages.User, error) {
	u := &storages.User{ID: getUserID(ctx)}
	if err := db.Db.Model(u).Select(); err != nil {
		return nil, err
	}
	return u, nil
}
