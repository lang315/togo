package transport

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lang315/togo/internal/storages"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type AuthRouter struct {
	app     *echo.Echo
	db      *storages.Data
	jwtAuth echo.MiddlewareFunc
}

func (a *AuthRouter) Install(app *echo.Echo, db *storages.Data, i... interface{}) {
	a.app = app
	a.db = db
	a.jwtAuth = middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserClaims{},
		SigningKey: []byte(jwtSecret),
	})

	r := a.app.Group("/auth")
	r.POST("/login", a.login)
	r.GET("/refresh", a.refresh, a.jwtAuth)
}

func (a *AuthRouter) refresh(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*UserClaims)
	token, err := makeUserJWTToken(claims.UserID)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, echo.Map{"token": token})
}

func (a *AuthRouter) login(ctx echo.Context) error {
	f := &struct {
		ID       string `json:"id"`
		Password string `json:"password"`
	}{}
	if err := ctx.Bind(f); err != nil {
		return err
	}
	f.ID = strings.ToLower(f.ID)

	u := &storages.User{}
	if err := a.db.Db.Model(u).Where("id = ?", f.ID).Select(u); err != nil {
		return echo.ErrUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(f.Password)); err != nil {
		return echo.ErrUnauthorized
	}

	token, err := makeUserJWTToken(u.ID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, echo.Map{"token": token})
}
