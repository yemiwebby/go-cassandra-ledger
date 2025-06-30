CREATE KEYSPACE IF NOT EXISTS ledger
WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE ledger;

CREATE TABLE ledger_entries (
    legal_entity text,
    namespace text,
    name text,
    currency text,
    account_id text,
    time_bucket text,
    committed_ts timestamp,
    type text,
    amount double,
    description text,
    PRIMARY KEY ((legal_entity, namespace, name, currency, account_id, time_bucket), committed_ts)
);
