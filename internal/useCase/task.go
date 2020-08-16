package useCase

import (
	"fmt"
	"github.com/lang315/togo/internal/bolt"
	"strconv"
	"sync"
	"time"
)

const layoutISO = "2006-01-02"

type Tasker struct {
	db                 *bolt.Bolt
	locker             *sync.RWMutex
	bucketName         string
	CountMax, CountMin int
}

func NewTasker(db *bolt.Bolt, countMax int) *Tasker {
	return &Tasker{
		db:       db,
		locker:   new(sync.RWMutex),
		CountMax: countMax,
		CountMin: 0,
	}
}

func (t *Tasker) AddNewID(id string) error {
	return t.db.Update(id, 1)
}

func (t *Tasker) Count(id string) int {
	count := 0
	k, err := t.db.GetValue(id)
	if err != nil {
		return -1
	}

	count, _ = strconv.Atoi(string(k))

	return count
}

func (t *Tasker) Update(id string, count int) error {
	return t.db.Update(id, count)
}

func (t *Tasker) ResetLimit() {

	deadline := t.timeReset()
	for {
		time.Sleep(deadline)
		t.reset()
		deadline = t.timeReset()
	}
}

func (t *Tasker) reset() {
	listID := make([][]byte, 0)
	listID = t.db.GetAllValues()

	if len(listID) < 1 {
		return
	}

	for i := 0; i < len(listID); i++ {
		t.db.Update(string(listID[i]), 0)
	}
}

func (t *Tasker) timeReset() time.Duration {
	now := time.Now()
	y, m, d := now.Add(time.Hour * 24).Date()
	deadlineString := fmt.Sprintf("%d-%s-%dT0:0:0", y, m, d)
	deadline, _ := time.Parse(time.RFC3339, deadlineString)
	result := time.Since(deadline)
	return result
}
