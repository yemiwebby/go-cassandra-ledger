CREATE KEYSPACE IF NOT EXISTS ledger
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE ledger;

CREATE TABLE IF NOT EXISTS ledger_entries (
    account_id TEXT,
    time_bucket TEXT,
    committed_ts TIMESTAMP,
    type TEXT,
    amount DECIMAL,
    description TEXT,
    PRIMARY KEY ((account_id, time_bucket), committed_ts)
) WITH CLUSTERING ORDER BY (committed_ts DESC);
