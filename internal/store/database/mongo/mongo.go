package mongo

import (
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
)

type MongoStore struct {
	db string
}

func NewMongoStore() (*MongoStore, error) {

	return &MongoStore{db: ""}, nil
}

func (c *MongoStore) WriteLedgerEntry(entry entry.LedgerEntry) error {
	return nil
}
