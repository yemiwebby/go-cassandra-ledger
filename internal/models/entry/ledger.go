package entry

import (
	"errors"
	"math"

	"github.com/yemiwebby/go-cassandra-ledger/internal/models/config"
	"github.com/yemiwebby/go-cassandra-ledger/internal/models/input"
)

type LedgerEntry struct {
	Address       config.LedgerAddress
	Type          string // credit or debit
	Amount        float64
	Description   string
	Timestamp     int64
	ReportingTime *int64
}

type EntrySet struct {
	Entries []LedgerEntry
}

var (
	ErrUnbalancedEntrySet = errors.New("entry set is not balanced")
	ErrInvalidEntryType   = errors.New("entry set contains invalid entry type")
)

const (
	TypeCredit = "credit"
	TypeDebit  = "debit"
)

func NewLedgerEntry(input input.EntrySetInput) EntrySet {
	var entries []LedgerEntry
	for _, entry := range input.Entries {
		rTime := int64(0)
		if entry.ReportingTime != nil {
			rTime = *entry.ReportingTime
		}

		entries = append(entries, LedgerEntry{
			Address: config.LedgerAddress{
				LegalEntity: entry.Address.LegalEntity,
				Namespace:   entry.Address.Namespace,
				Name:        entry.Address.Name,
				Currency:    entry.Address.Currency,
				AccountID:   entry.Address.AccountID,
			},
			Type:          entry.Type,
			Amount:        entry.Amount,
			Description:   entry.Description,
			Timestamp:     entry.Timestamp,
			ReportingTime: &rTime,
		})
	}

	return EntrySet{Entries: entries}
}

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
