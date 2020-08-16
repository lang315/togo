package useCase

import (
	"github.com/lang315/togo/internal/bolt"
	"os"
	"testing"
)

var (
	bucketName = "task"
)

func initBolt() (*bolt.Bolt, error) {
	return bolt.NewBolt("testDB.txt", "taks")
}

func TestTask(t *testing.T) {
	db, err := initBolt()
	if err != nil {
		t.Error(err)
	}

	defer os.Remove("data.db")
	defer db.Close()

	id := "langdethuong"
	task := NewTasker(db, 5)
	err = task.AddNewID(id)
	if err != nil {
		t.Error(err)
	}

	count := task.Count(id)
	if count == -1 {
		t.Errorf("count failed")
	}

	countNew := count + 1
	err = task.Update(id, countNew)
	if err != nil {
		t.Error(err)
	}

	count = task.Count(id)
	if count != countNew {
		t.Errorf("upate failed")
	}
}
