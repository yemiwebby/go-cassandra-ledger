package entry_test

import (
	"testing"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
)

func TestEntrySet_IsBalanced(t *testing.T) {
	es := entry.EntrySet{
		Entries: []entry.LedgerEntry{
			{Amount: 100.0, Type: entry.TypeCredit},
			{Amount: 100.0, Type: entry.TypeDebit},
		},
	}

	ok, err := es.IsBalanced()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Errorf("expected entry set to be balanced")
	}
}
