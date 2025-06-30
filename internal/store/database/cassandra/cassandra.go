package cassandra

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/entry"
)

type CassandraStore struct {
	session *gocql.Session
}

func NewCassandraStore() (*CassandraStore, error) {
	cluster := gocql.NewCluster("localhost")
	cluster.Keyspace = "ledger"
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.RetryPolicy = &gocql.SimpleRetryPolicy{NumRetries: 3}

	session, err := cluster.CreateSession()
	if err == nil {
		log.Println("Connected to Cassandra")
		return &CassandraStore{session: session}, nil
	}

	return nil, fmt.Errorf("failed to connect to Cassandra after retries: %w", err)
}

func (c *CassandraStore) WriteLedgerEntry(entry entry.LedgerEntry) error {
	timeBucket := time.UnixMilli(entry.Timestamp).Format("2006-01")

	query := `INSERT INTO ledger_entries (
		legal_entity, namespace, name, currency, account_id,
		time_bucket, committed_ts, type, amount, description
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	return c.session.Query(query,
		entry.Address.LegalEntity,
		entry.Address.Namespace,
		entry.Address.Name,
		entry.Address.Currency,
		entry.Address.AccountID,
		timeBucket,
		time.UnixMilli(entry.Timestamp),
		entry.Type,
		entry.Amount,
		entry.Description,
	).Exec()
}

func (c *CassandraStore) GetEntries(
	addr config.LedgerAddress,
	timeAxis string,
	start, end time.Time,
) ([]entry.LedgerEntry, error) {

	timeLayout := "2006-01"
	startBucket := start.Format(timeLayout)
	endBucket := end.Format(timeLayout)

	query := `SELECT account_id, type, amount, description, committed_ts, reporting_ts
			  FROM ledger_entries
			  WHERE account_id = ? AND time_bucket >= ? AND time_bucket <= ?`

	iter := c.session.Query(query, addr.AccountID, startBucket, endBucket).Iter()

	var results []entry.LedgerEntry
	var accountID, typ, desc string
	var amount float64
	var committedTs, reportingTs time.Time

	for iter.Scan(&accountID, &typ, &amount, &desc, &committedTs, &reportingTs) {
		timestamp := committedTs.UnixMilli()
		if timeAxis == "reporting" {
			timestamp = reportingTs.UnixMilli()
		}

		// Filter by the actual timestamp range
		ts := time.UnixMilli(timestamp)
		if ts.Before(start) || ts.After(end) {
			continue
		}

		rTime := reportingTs.UnixMilli()
		results = append(results, entry.LedgerEntry{
			Address:       addr,
			Type:          typ,
			Amount:        amount,
			Description:   desc,
			Timestamp:     committedTs.UnixMilli(),
			ReportingTime: &rTime,
		})
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return results, nil
}
