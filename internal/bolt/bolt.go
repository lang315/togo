package bolt

import (
	"github.com/lang315/togo/internal/storages"
	"github.com/pkg/errors"
	"go.etcd.io/bbolt"
	"strconv"
)

type Bolt struct {
	db         *bbolt.DB
	BucketName string
}

func NewBolt(path, bucketName string) (*Bolt, error) {
	db, err := bbolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Bolt{
		db:         db,
		BucketName: bucketName,
	}, nil
}

func (b *Bolt) Migrate(data *storages.Data) error {
	listUser, err := data.AllUsers()
	if err != nil {
		return err
	}

	for i := 0; i < len(listUser); i++ {
		b.Update(listUser[i].ID, 0)
	}

	return nil
}

func (b *Bolt) Update(id string, count int) error {
	return b.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.BucketName))
		return b.Put([]byte(id), []byte(strconv.Itoa(count)))
	})
}

func (b *Bolt) GetValue(key string) ([]byte, error) {

	v := make([]byte, 0)
	err := b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.BucketName))
		v = b.Get([]byte(key))
		if v == nil {
			return errors.New("key is not existing")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return v, nil
}

func (b *Bolt) GetAllValues() [][]byte {
	listValue := make([][]byte, 0)
	b.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(b.BucketName))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; {
			listValue = append(listValue, k)
		}
		return nil
	})

	return listValue
}

func (b *Bolt) Close() {
	b.db.Close()
}
