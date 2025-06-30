# A backend focused ledger microservice built with Go and Cassandra.

This project was built to explore how FinTech companies like Monzo build scalable, auditable, and append-only ledger systems using Go and Cassandra — and to serve as a practical, testable foundation for financial backend infrastructure.

## Objectives

This mini project aims to:

- Simulate a high-volume, append-only financial ledger
- Demonstrate use of Cassandra for scalable, distributed time-series data

## What is a Ledger?

A ledger is an immutable, append-only data store used to track financial transactions (Credits and Debits) or any state change over time. It is a foundational component in many financial systems, providing a reliable and auditable record of all transactions.

## Ledger Microservice Core Concepts

- **Immutable**: Once a transaction is recorded, it cannot be altered or deleted. Append-only. No updates, just new entries.
- **Append-Only**: New transactions are added to the end of the ledger.
- **Time-series Friendly**: Designed to handle time-series data efficiently, making it suitable for tracking changes over time. Sorted by timestamp.
- - Balance calculation is derived from entries. Currently, only full-scan computation is implemented; precomputed blocks and snapshots are future improvements. The ledger does not store the current balance, but rather the individual transactions that can be used to calculate it.
- **Event Sourcing**: The ledger can be used as an event store, where each transaction represents an event in the system. This allows for replaying events to reconstruct the state of the system at any point in time.

## Why Cassandra?

Cassandra is ideal for this because:

- High write throughput: It can handle a large number of writes per second, making it suitable for high-volume transaction systems. Great for appending lots of transactions.
- Horizontal scalability: It can scale out by adding more nodes, allowing it to handle large datasets and high traffic loads.
- Time-series data support: Cassandra's data model is well-suited for time-series data, making it easy to store and query transactions based on timestamps. Can partition by user and time bucket.
- Tunable consistency: You control read/write trade-offs.
- No Single point of failure: Good for always-on financial systems.

## Other Details

This project follows a clean, layered architecture. The HTTP layer is thin and declarative, delegating all domain logic to an Engine (Service Layer), which in turn interacts with a pluggable Store layer. This separation allows for testability, future extensibility, and Monzo-style production readiness.

## Database Connection

```sql
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

```

```sql
docker compose -f docker/cassandra-compose.yml down -v
docker compose -f docker/cassandra-compose.yml up -d --force-recreate
```

```
docker exec -it cassandra cqlsh
```

```sql
DESCRIBE KEYSPACES;
```

## Testing

Create a transaction.json file in the root of the project directory and paste the following content:

```
touch transaction.json
```

```json
{
  "entries": [
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.account",
        "name": "main",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 100,
      "type": "debit",
      "description": "Initial funding",
      "timestamp": 1727347200000
    },
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.equity",
        "name": "owners",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 100,
      "type": "credit",
      "description": "Owner's equity",
      "timestamp": 1727347200000
    },
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.account",
        "name": "main",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 50,
      "type": "debit",
      "description": "Client payment",
      "timestamp": 1727520000000
    },
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.revenue",
        "name": "general",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 50,
      "type": "credit",
      "description": "Income from service",
      "timestamp": 1727520000000
    },
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.utilities",
        "name": "electric",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 30,
      "type": "debit",
      "description": "Electricity bill",
      "timestamp": 1727606400000
    },
    {
      "address": {
        "legal_entity": "fintech_uk",
        "namespace": "com.fintech.account",
        "name": "main",
        "currency": "GBP",
        "account_id": "*"
      },
      "amount": 30,
      "type": "credit",
      "description": "Paid utilities",
      "timestamp": 1727606400000
    }
  ]
}
```

```bash
curl -X POST http://localhost:8080/transactions \
  -H "Content-Type: application/json" \
  -d @transaction.json
```

### Check the Ledger Entries

```bash
docker exec -it cassandra cqlsh
```

```bash
USE ledger;
SELECT * FROM ledger_entries;
```

```sql
cqlsh:ledger> SELECT \* FROM ledger_entries;

 legal_entity | namespace             | name     | currency | account_id | time_bucket | committed_ts                    | amount | description         | type
--------------+-----------------------+----------+----------+------------+-------------+---------------------------------+--------+---------------------+--------
   fintech_uk |    com.fintech.equity |   owners |      GBP |          * |     2024-09 | 2024-09-26 10:40:00.000000+0000 |    100 |      Owner's equity | credit
   fintech_uk | com.fintech.utilities | electric |      GBP |          * |     2024-09 | 2024-09-29 10:40:00.000000+0000 |     30 |    Electricity bill |  debit
   fintech_uk |   com.fintech.revenue |  general |      GBP |          * |     2024-09 | 2024-09-28 10:40:00.000000+0000 |     50 | Income from service | credit
   fintech_uk |   com.fintech.account |     main |      GBP |          * |     2024-09 | 2024-09-26 10:40:00.000000+0000 |    100 |     Initial funding |  debit
   fintech_uk |   com.fintech.account |     main |      GBP |          * |     2024-09 | 2024-09-28 10:40:00.000000+0000 |     50 |      Client payment |  debit
   fintech_uk |   com.fintech.account |     main |      GBP |          * |     2024-09 | 2024-09-29 10:40:00.000000+0000 |     30 |      Paid utilities | credit
```

(6 rows)

### Truncate the ledger_entries table

```sql
TRUNCATE ledger_entries;
```

## Primary Key Structure in ledger_entries

The table uses a composite primary key:

```sql
PRIMARY KEY ((legal_entity, namespace, name, currency, account_id, time_bucket), committed_ts)
```

This breaks down as:

- Partition Key: (account_id, time_bucket)

  - Groups entries by account and by monthly bucket (e.g., 2024-09)

  - Ensures efficient writes and distributes data across the cluster

- Clustering Column: committed_ts

  - Orders entries within each partition chronologically

  - Allows efficient time-range queries within a specific account/month

This structure enables append-only, time-series storage per account while keeping the write and read paths performant.

### EntrySet Balancing Logic

In a double-entry ledger, each transaction must be balanced:

The total amount of credits must equal the total amount of debits.

The service enforces this rule via the IsBalanced() method on the EntrySet:

```go
func (es EntrySet) IsBalanced() (bool, error) {
var total float64

    for _, entry := range es.Entries {
        switch entry.Type {
        case TypeCredit:
            total += entry.Amount
        case TypeDebit:
            total -= entry.Amount
        default:
            return false, ErrInvalidEntryType
        }
    }

    if math.Abs(total) > 0.00001 {
        return false, ErrUnbalancedEntrySet
    }

    return true, nil

}
```

- Credits increase the total

- Debits decrease the total

- A transaction is considered balanced if the final total is close to zero (within a small epsilon margin for floating-point rounding)

- If unbalanced, the request is rejected with an error

This logic ensures the integrity of the ledger and prevents invalid financial entries.

## Running the Application

### Run Cassandra in Docker

If you don't already have a running Cassandra container:

```bash
docker run --name cassandra-ledger \
  -p 9042:9042 \
  -d cassandra:4.1

```

This exposes port 9042, which is what your Go app will connect to.

### Wait for Cassandra to be Ready

Cassandra takes 20-40s to fully boot up. You can check logs:

```bash
docker logs -f cassandra-ledger
```

Look for:

```bash
Startup complete
```

### Check if Keyspace & Table Exist

Use `cqlsh`:

```
docker exec -it cassandra-ledger cqlsh
```

Inside the shell:

```sql
DESCRIBE KEYSPACES;

USE ledger;

DESCRIBE TABLES;

SELECT * FROM ledger_entries LIMIT 5;
```

### Run the Go Project

```bash
go run main.go
```

### Test the /transaction Endpoint

Use this as an example POST:

```bash
curl -X POST http://localhost:8080/transaction \
  -H "Content-Type: application/json" \
  -d @transaction.json
```

### Test the /balance Endpoint

```bash
curl "http://localhost:8080/balance?name=customer-facing-balance&start=2024-01-01T00:00:00Z&end=2025-12-31T23:59:59Z"

```

## Running Tests

Unit Tests

```bash
go test ./... -v
```

Integration Tests

```bash
go test ./test/integration -v
```

## Future Improvements

- Implement `balance_blocks` for fast partial-sum lookups
- Add support for historical snapshots (`as_of` balance)
- Kafka-based ingestion for event-sourced pipelines
- Enhanced testing with time-skewed entries and edge cases
- Parallelized reads for high-volume accounts
