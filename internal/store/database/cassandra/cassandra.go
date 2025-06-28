package cassandra

import (
	"log"
	"time"

	gocql "github.com/apache/cassandra-gocql-driver/v2"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
)

type CassandraStore struct {
	session *gocql.Session
}

func NewCassandraStore() (*CassandraStore, error) {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "ledger"
	cluster.Consistency = gocql.Quorum

	var err error
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
	}

	log.Println("Connected to Cassandra")

	return &CassandraStore{session: session}, nil
}

func (c *CassandraStore) WriteLedgerEntry(entry entry.LedgerEntry) error {
	timeBucket := time.UnixMilli(entry.Timestamp).Format("2006-01")

	query := `INSERT INTO ledger_entries (
	   account_id, time_bucket, commited_ts, type, amount, description
	) VALUES (?, ?, ?, ?, ?, ?)`

	return c.session.Query(query,
		entry.AccountID,
		timeBucket,
		time.UnixMilli(entry.Timestamp),
		entry.Type,
		entry.Amount,
		entry.Description,
	).Exec()
}
