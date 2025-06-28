package store

import "github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"

type LedgerEntryStore interface {
	WriteLedgerEntry(entry entry.LedgerEntry) error
}
