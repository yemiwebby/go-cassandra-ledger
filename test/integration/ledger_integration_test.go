package integration

import (
	"testing"
	"time"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
	"github.com/yemiwebby/go-cassandra-ledger/internal/service"
	"github.com/yemiwebby/go-cassandra-ledger/internal/store/database/cassandra"
)

func TestIntegration_InsertAndReadLedgerEntry(t *testing.T) {
	store, err := cassandra.NewCassandraStore()
	if err != nil {
		t.Fatalf("failed to connect to Cassandra: %v", err)
	}
	defer store.Close()

	engine := service.NewEngine(store)

	now := time.Now().UTC()
	entrySet := entry.EntrySet{
		Entries: []entry.LedgerEntry{
			{
				Address: config.LedgerAddress{
					AccountID:   "integration-test-users",
					LegalEntity: "fintech_uk",
					Namespace:   "com.test.integration",
					Name:        "test_account",
					Currency:    "GBP",
				},
				Amount:      42.0,
				Type:        entry.TypeCredit,
				Description: "Integration test credit",
				Timestamp:   now.UnixMilli(),
			},
		},
	}

	err = engine.ProcessEntrySet(entrySet)
	if err != nil {
		t.Fatalf("failed to process entry set: %v", err)
	}
}
