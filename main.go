package main

import (
	"context"
	"github.com/go-pg/pg/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lang315/togo/internal/bolt"
	"github.com/lang315/togo/internal/storages"
	"github.com/lang315/togo/internal/transport"
)

func main() {

	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "Conghuy.315@",
		Database: "togo",
	})

	if err := db.Ping(context.TODO()); err != nil {
		panic(err)
	}
	data := &storages.Data{Db: db}
	data.Migrate()

	boltDB, err := bolt.NewBolt("rateLimit.db", "task")
	if err != nil {
		panic(err)
	}

	err = boltDB.Migrate(data)
	if err != nil {
		panic(err)
	}

	installRouters := []transport.BaseRouter{
		&transport.AuthRouter{},
		&transport.TaskRouter{},
	}

	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	for i := 0; i < len(installRouters); i++ {
		installRouters[i].Install(app, data, boltDB)
	}

	app.Logger.Fatal(app.Start(":1323"))
}
