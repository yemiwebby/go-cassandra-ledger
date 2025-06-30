## Overview

`go-cassandra-ledger` is a Monzo-inspired, event-driven financial ledger system built in Go (Golang) and backed by Cassandra. It implements a production-grade, double-entry, append-only ledger with time-bucketed storage, balance definition abstraction, and a scalable read/write pattern for modern FinTech systems.

This project closely models the real architecture published by Monzo in their [Engineering Blog](https://monzo.com/blog/2023/04/28/speeding-up-our-balance-read-time-the-planning-phase), and serves as an open-source educational tool with production-ready constraints, patterns, and solutions considered part of the MVP.

## Design Goals

- Model a robust double-entry ledger with strong transactional integrity.
- Achieve high write throughput via an append-only Cassandra schema.
- Support multiple balance types (e.g. `customer-facing`, `interest-chargeable`).
- Enable fast balance computation using precomputed blocks and time bucketing.
- Allow time-based filtering using `committed` or `reporting` axes.
- Embrace real-world production patterns from Monzo’s architecture.
- Maintain auditability with immutable entries and traceable transactions.

## Key Concepts & Terminology

- **Ledger**: Source of truth for all monetary transactions, using append-only double-entry bookkeeping.
- **EntrySet**: A group of ledger entries representing a complete money movement.
- **Ledger Address**: A unique 5-tuple key (`legal_entity`, `namespace`, `name`, `currency`, `account_id`) identifying an account.
- **Balance Name**: Logical label representing a computed balance (e.g. `customer-facing-balance`).
- **Committed Timestamp**: When the entry is persisted; used for partitioning.
- **Reporting Timestamp**: When the transaction takes effect for accounting purposes.
- **Time Axis**: Determines which timestamp to use when querying balances.
- **Balance Block**: Precomputed aggregates per account and time bucket to speed up reads.

## Architecture Overview

- **Ingestion API**: Accepts transactions (`EntrySets`) and persists them to Cassandra.
- **Balance Engine**: Computes balances via full scans or block-based optimizations.
- **Config Loader**: Loads `address_config.yml` and `balance_definitions.yml`.
- **Ledger Query API**: Exposes `/transactions`, `/balance`, and `/health` endpoints.
- **Healthcheck Endpoint**: `/health` route for liveness and monitoring probes.

## Data Models

1. **Ledger Entries Table (`ledger_entries`)**

| Field        | Type                               |
| ------------ | ---------------------------------- |
| account_id   | TEXT                               |
| time_bucket  | TEXT (e.g. 2025-06)                |
| committed_ts | TIMESTAMP                          |
| reporting_ts | TIMESTAMP                          |
| txn_id       | UUID                               |
| type         | TEXT (credit or debit)             |
| amount       | DECIMAL                            |
| address      | TEXT (flattened string of 5-tuple) |
| description  | TEXT                               |

```sql
PRIMARY KEY ((account_id, time_bucket), committed_ts)
```

<!--
2. **Balance Blocks Table (`balance_blocks`)**

| Field        | Type      |
| ------------ | --------- |
| account_id   | TEXT      |
| time_bucket  | TEXT      |
| balance_name | TEXT      |
| partial_sum  | DECIMAL   |
| last_txn_ts  | TIMESTAMP |

Purpose:

- Speed up balance reads by skipping full historical scans.
- Blocks can be merged during snapshot creation.

3. **Balance Snapshots Table (`balance_snapshots`)**

| Field            | Type      |
| ---------------- | --------- |
| account_id       | TEXT      |
| balance_name     | TEXT      |
| snapshot_date    | DATE      |
| snapshot_balance | DECIMAL   |
| computed_at      | TIMESTAMP |

Purpose:

- Supports queries like: "What was my balance on Jan 1, 2024?" -->

4. Config Files (Versioned)
   `address_config.yml`

```yaml
main_account_gbp:
  legal_entity: fintech_uk
  namespace: com.fintech.account
  name: main
  currency: GBP
  account_id: "*"

revenue_account_gbp:
  legal_entity: fintech_uk
  namespace: com.fintech.revenue
  name: general
  currency: GBP
  account_id: "*"

utilities_account_gbp:
  legal_entity: fintech_uk
  namespace: com.fintech.utilities
  name: electric
  currency: GBP
  account_id: "*"
```

`balance_definitions.yml`

```yaml
customer-facing-balance:
  time_axis: committed
  addresses:
    - main_account_gbp

interest-chargeable-balance:
  time_axis: committed
  addresses:
    - revenue_account_gbp

utilities-expense-balance:
  time_axis: committed
  addresses:
    - utilities_account_gbp
```

## API Endpoints

| Method | Endpoint                            | Description                             |
| ------ | ----------------------------------- | --------------------------------------- |
| POST   | `/transaction`                      | Ingest an `EntrySet` (must be balanced) |
| GET    | `/balance?name=X&start=...&end=...` | Compute balance by logical name         |
| GET    | `/health`                           | Health check (returns 200 OK if alive)  |

## Balance Computation Logic

1. Full Scan (MVP)

- Fetch all entries matching balance address set.

- Filter entries based on committed_ts or reporting_ts.

- Sum based on credit/debit type.

2. (Future) Block-Based

- Use balance_blocks to sum historical buckets.

- Read most recent entries from ledger_entries for delta.

- Merge results for fast reads.

3. (Future) Snapshot + Delta

- Store daily/monthly snapshots.

- Use as base + apply delta from recent entries.

## Partitioning & Bucketing

| Strategy      | Value                                |
| ------------- | ------------------------------------ |
| Partition Key | `account_id`, `time_bucket`          |
| Clustering    | `committed_ts DESC`                  |
| Bucket Size   | 1 month                              |
| Benefit       | Reduces partition size, speeds reads |

## Handling External Money Movement: Inbound Transfers

In real-world banking systems, not all transactions originate within the institution's ledger. For example, a Monzo customer may receive money from an external bank or card network. Even in these cases, double-entry integrity must be preserved.

To support this, we introduce the concept of synthetic internal ledger addresses to represent external systems.

### Example: External Transfer into Customer Account

| Type   | Account ID            | Amount | Description                  |
| ------ | --------------------- | ------ | ---------------------------- |
| Debit  | external_inbound_bank | £100   | Inbound from external source |
| Credit | fintech_uk_main       | £100   | Received from external party |

This EntrySet models a complete money movement:

```json
{
  "entry_set_id": "uuid-7890",
  "entries": [
    {
      "account_id": "external_inbound_bank",
      "type": "debit",
      "amount": 100.0,
      "description": "Inbound transfer from external bank"
    },
    {
      "account_id": "bob456",
      "type": "credit",
      "amount": 100.0,
      "description": "Received from Alice"
    }
  ]
}
```

### Synthetic Ledger Address: External Inbound Bank

To support this logic, the system includes predefined internal ledger addresses for external systems:

```yaml
external_inbound_bank:
  legal_entity: fintech_uk
  namespace: com.fintech.inbound
  name: settlement
  currency: GBP
  account_id: external
```

This ensures:

- All money movement remains traceable within the ledger.
- EntrySets are always balanced (sum to zero).
- External payments rails (e.g FPS, Mastercard, BACS) can be reconciled against this account.

### Benefits

- Auditability: Every penny is accounted for, even across system boundaries.
- Consistency: Double-entry integrity is maintained system-wide.
- Extendability: Enables modeling of external-to-internal and internal-to-external transactions symmetrically.

## Testing & Simulation

| Phase   | Description                                |
| ------- | ------------------------------------------ |
| Phase 1 | Manual curl + simulated `transaction.json` |
| Phase 2 | Load tester script (bash or Go-based)      |
| Phase 3 | Kafka simulator (event ingestion post-MVP) |

### Phase 1: Bash Script

- Simulate realistic traffic by firing 1000+ transactions per user.
- Vary between credit and debit, random timestamps.
- Validate ingestion and read paths.

### Phase 2: Kafka Ingestion Simulator (Post-MVP)

- Build event-based simulator in Go.
- Publishes EntrySet events (salary, ATM withdrawal, refunds, etc.)
- `go-cassandra-ledger` listens via Kafka consumer and ingests live.

## High-Level Diagram

This diagram illustrates the high-level architecture of the `go-cassandra-ledger` system, showing the key components and their interactions:
![High-Level Architecture](./docs/images/high-level.png)

## Future Improvements

| Feature                   | Benefit                                         |
| ------------------------- | ----------------------------------------------- |
| **Balance Blocks Table**  | Enables partial precomputed sums for fast reads |
| **Snapshots Table**       | Allows historical `as_of` balance queries       |
| Kafka consumer            | Ingest ledger entries from external systems     |
| Real-time balance stream  | Push-based update events for balance watchers   |
| Fraud rule scoring        | Detect anomalies at ingest time                 |
| Parallel reads per bucket | Boost performance for high-volume users         |
| Custom balance views      | Enable dynamic grouping/filtering               |
