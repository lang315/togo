package transport

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lang315/togo/internal/bolt"
	"github.com/lang315/togo/internal/storages"
	"github.com/lang315/togo/internal/useCase"
	"net/http"
)

type TaskRouter struct {
	app     *echo.Echo
	db      *storages.Data
	jwtAuth echo.MiddlewareFunc
	tasker  *useCase.Tasker
	boltDB  *bolt.Bolt
}

func (t *TaskRouter) Install(app *echo.Echo, db *storages.Data, i ...interface{}) {
	t.app = app
	t.db = db
	t.jwtAuth = middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &UserClaims{},
		SigningKey: []byte(jwtSecret),
	})

	if len(i) > 0 {
		t.boltDB = i[0].(*bolt.Bolt)
	}

	t.tasker = useCase.NewTasker(t.boltDB, 5)

	r := t.app.Group("/task", t.jwtAuth)
	r.POST("", t.add)
}

func (t *TaskRouter) add(ctx echo.Context) error {
	task := &storages.Task{}

	if err := ctx.Bind(task); err != nil {
		return err
	}

	userID := getUserID(ctx)
	ok, user := t.db.ValidateUser(userID)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, "")
	}

	countTask := t.tasker.Count(userID)
	if countTask == -1 {
		t.tasker.AddNewID(userID)
		return ctx.JSON(http.StatusOK, "Done")
	}

	if countTask > user.MaxToDo {
		return ctx.JSON(http.StatusTooManyRequests, "Too Many Requests")
	}

	countTask++
	t.tasker.Update(userID, countTask)

	return ctx.JSON(http.StatusOK, "Done")
}
