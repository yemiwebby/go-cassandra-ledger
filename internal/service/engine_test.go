package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
	"github.com/yemiwebby/go-cassandra-ledger/internal/service"
)

type mockStore struct {
	shouldFail bool
	entries    []entry.LedgerEntry
}

func (m *mockStore) WriteLedgerEntry(e entry.LedgerEntry) error {
	if m.shouldFail {
		return errors.New("mock write failure")
	}
	m.entries = append(m.entries, e)
	return nil
}

func (m *mockStore) GetEntries(addr config.LedgerAddress, timeAxis string, start, end time.Time) ([]entry.LedgerEntry, error) {
	return nil, nil
}

func TestEngine_ProcessEntrySet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockStore{}
		engine := service.NewEngine(mock)

		es := entry.EntrySet{
			Entries: []entry.LedgerEntry{
				{Amount: 100, Type: entry.TypeCredit},
				{Amount: 100, Type: entry.TypeDebit},
			},
		}

		err := engine.ProcessEntrySet(es)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(mock.entries) != 2 {
			t.Errorf("expected 2 entries written, got %d", len(mock.entries))
		}
	})

	t.Run("failure", func(t *testing.T) {
		mock := &mockStore{shouldFail: true}
		engine := service.NewEngine(mock)

		es := entry.EntrySet{
			Entries: []entry.LedgerEntry{
				{Amount: 100, Type: entry.TypeCredit},
			},
		}

		err := engine.ProcessEntrySet(es)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
