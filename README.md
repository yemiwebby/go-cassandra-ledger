# A backend focused ledger microservice built with Go and Cassandra.
This is a lightweight, financial ledger service built with Go and Cassandra. It models an immutable, append-only transaction system designed for high-write scalability and time-series access. The project demonstrates clean backend architecture, event-style ledger recording, and efficient balance computation, all optimized for modern Fintech systems.

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
- Balance calculation is derived, not stored. The ledger does not store the current balance, but rather the individual transactions that can be used to calculate it.
- **Event Sourcing**: The ledger can be used as an event store, where each transaction represents an event in the system. This allows for replaying events to reconstruct the state of the system at any point in time.

## Why Cassandra?

Cassandra is ideal for this because:
- High write throughput: It can handle a large number of writes per second, making it suitable for high-volume transaction systems. Great for appending lots of transactions.
- Horizontal scalability: It can scale out by adding more nodes, allowing it to handle large datasets and high traffic loads.
- Time-series data support: Cassandra's data model is well-suited for time-series data, making it easy to store and query transactions based on timestamps. Can partition by user and time bucket.
- Tunable consistency: You control read/write trade-offs.
- No Single point of failure: Good for always-on financial systems.