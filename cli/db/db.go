package db

import (
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type todo struct {
	Key   uint64
	Value string
}

type DB interface {
	Add(string) error
	List() ([]todo, error)
	Do(int64) error
}

type Boltdb struct {
	DB
	database    *bolt.DB
	initialized bool
	bucketName  string
}

type BoltOption func(*bolt.Options)

func WithTimeout(timeout time.Duration) BoltOption {
	return func(o *bolt.Options) {
		o.Timeout = timeout
	}
}

func (db *Boltdb) init(datafile string, opts []BoltOption) {
	bo := bolt.Options{}

	for _, o := range opts {
		o(&bo)
	}

	conn, err := bolt.Open(datafile, 0600, &bo)

	if err != nil {
		log.Fatalf("failed to open database:\n%+v", err)
	}

	err = conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(db.bucketName))

		return err
	})

	if err != nil {
		log.Fatalf("failed create bucket:\n%+v", err)
	}

	db.database = conn

	db.initialized = true
}

func (db *Boltdb) bucket(tx *bolt.Tx) *bolt.Bucket {
	return tx.Bucket([]byte(db.bucketName))
}

func (db *Boltdb) Close() {
	db.database.Close()
}

func New(datafile string, opts ...BoltOption) *Boltdb {
	db := Boltdb{bucketName: "Todos"}

	db.init(datafile, opts)

	return &db
}

func (db *Boltdb) Add(item string) error {
	if !db.initialized {
		return bolt.ErrDatabaseNotOpen
	}

	return db.database.Update(func(tx *bolt.Tx) error {
		b := db.bucket(tx)

		key, _ := b.NextSequence()

		formattedKey := itob(key)

		fmt.Printf("Key will be %v\n", key)

		return b.Put(formattedKey, []byte(item))
	})
}

func (db *Boltdb) List() ([]todo, error) {
	if !db.initialized {
		return nil, bolt.ErrDatabaseNotOpen
	}

	var result []todo

	err := db.database.View(func(tx *bolt.Tx) error {
		b := db.bucket(tx)

		return b.ForEach(func(k, v []byte) error {
			result = append(result, todo{btoi(k), string(v)})
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (db *Boltdb) Do(id int64) error {
	if !db.initialized {
		return bolt.ErrDatabaseNotOpen
	}

	return db.database.Update(func(tx *bolt.Tx) error {
		b := db.bucket(tx)

		return b.Delete(itob(uint64(id)))
	})
}

func itob(i uint64) []byte {
	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, i)

	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}
