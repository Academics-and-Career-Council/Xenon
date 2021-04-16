package internal

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type badgerClient struct {
	d *badger.DB
}

var BadgerDB badgerClient

func Init() {
	db, err := badger.Open(badger.DefaultOptions("./cache"))
	if err != nil {
		log.Fatal(err)
	}
	// b.db = db
	// err = b.db.View(func(txn *badger.Txn) error {
	// 	_, err := txn.Get([]byte("hello"))
	// 	return err
	// })
	BadgerDB = badgerClient{db}
}

func (b badgerClient) Save(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = b.d.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(p)).WithTTL(time.Hour * 24)
		err := txn.SetEntry(e)
		return err
	})
	return err
}

func (b badgerClient) Get(key string, dest interface{}) error {
	var p []byte

	err := b.d.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))

		if err != nil {
			log.Print(err)
			return err
		}

		// err = item.Value(func(val []byte) error {
		// 	p = append([]byte{}, val...)
		// 	return err
		// })

		return nil
	})
	if err != nil {
		return err
	}
	err = json.Unmarshal(p, dest)
	return err
}
