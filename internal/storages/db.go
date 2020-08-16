package storages

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"log"
)

type Data struct {
	Db *pg.DB
}

func (d *Data) Migrate() {
	d.createTableIfNotExists(&User{})
	d.createTableIfNotExists(&Task{})
}

func (d *Data) createTableIfNotExists(tbl interface{}) {
	if err := d.Db.Model(tbl).CreateTable(&orm.CreateTableOptions{IfNotExists: true}); err != nil {
		log.Fatalln(err)
	}
}

func (d *Data) ValidateUser(userID string) (bool, *User) {
	user := &User{ID: userID}
	err := d.Db.Model(user).Select()
	if err != nil {
		return false, nil
	}

	return true, user
}

func (d *Data) AllUsers() ([]*User, error) {
	var users []*User
	err := d.Db.Model(&users).Select()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (d *Data) AddTask(task *Task) error {
	_, err := d.Db.Model(task).Insert()
	return err
}
