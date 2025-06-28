package service

import (
	"log"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
	"github.com/yemiwebby/go-cassandra-ledger/internal/store"
)

type Engine struct {
	Store store.LedgerEntryStore
}

func NewEngine(store store.LedgerEntryStore) *Engine {
	return &Engine{
		Store: store,
	}
}

func (e *Engine) ProcessEntrySet(set entry.EntrySet) error {
	for _, entry := range set.Entries {
		if err := e.Store.WriteLedgerEntry(entry); err != nil {
			log.Printf("Failed to write entry: %v", err)
			return err
		}
	}
	return nil
}
