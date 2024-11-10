package db

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type (
	DB struct {
		DB *leveldb.DB
	}
)

func NewDB() *DB {
	db, err := leveldb.OpenFile("blockchain-db", nil)
	if err != nil {
		log.Fatalf("failed to initialize db: %+v", err)
	}

	return &DB{
		DB: db,
	}
}

func (db *DB) Put(key []byte, value []byte) {
	db.DB.Put(key, value, nil)
}

func (db *DB) Get(key []byte) ([]byte, error) {
	return db.DB.Get(key, nil)
}
