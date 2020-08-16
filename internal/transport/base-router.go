package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/lang315/togo/internal/storages"
)

type BaseRouter interface {
	Install(app *echo.Echo, db *storages.Data, i... interface{})
}
