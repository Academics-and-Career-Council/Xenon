package internal

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dgraph-io/badger/v3"
)

type badgerClient struct {
	*badger.DB
}

var BadgerDB badgerClient

func (b badgerClient) Init() {
	db, err := badger.Open(badger.DefaultOptions("./cache"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	b.DB = db
}

func (b badgerClient) Save(key string, value interface{}) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = b.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(p)).WithTTL(time.Hour * 24)
		err := txn.SetEntry(e)
		return err
	})
	return err
}

func (b badgerClient) Get(key string, dest interface{}) error {
	var p []byte

	err := b.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("answer"))

		err = item.Value(func(val []byte) error {
			p = append([]byte{}, val...)
			return err
		})

		return err
	})
	if err != nil {
		return err
	}
	err = json.Unmarshal(p, dest)
	return err
}
