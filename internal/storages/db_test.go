package storages

import (
	"context"
	"github.com/go-pg/pg/v10"
	"testing"
)

func initDB() (*Data, error) {
	db := pg.Connect(&pg.Options{
		Addr:     ":5432",
		User:     "postgres",
		Password: "Conghuy.315@",
		Database: "togo",
	})

	if err := db.Ping(context.TODO()); err != nil {
		return nil, err
	}

	return &Data{Db: db}, nil
}

func TestData_AllUsers(t *testing.T) {
	db, err := initDB()
	if err != nil {
		t.Error(err)
	}

	db.Migrate()

	_, err = db.AllUsers()
	if err != nil {
		t.Error(err)
	}
}
