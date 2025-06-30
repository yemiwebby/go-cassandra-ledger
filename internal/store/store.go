package store

import (
	"time"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
)

type LedgerEntryStore interface {
	WriteLedgerEntry(entry entry.LedgerEntry) error
	GetEntries(addr config.LedgerAddress, timeAxis string, start, end time.Time) ([]entry.LedgerEntry, error)
}
