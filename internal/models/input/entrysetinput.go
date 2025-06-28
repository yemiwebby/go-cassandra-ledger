package input

type LedgerEntryInput struct {
	AccountID     string  `json:"account_id"`
	Type          string  `json:"type"` // credit/debit
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	Timestamp     int64   `json:"timestamp"`
	ReportingTime int64   `json:"reporting_time"`
}

/*
Example Entry
{
"entries": [
	{"account_id": "external", "type": "debit", "amount": 100},
	{"account_id": "user123", "type": "credit", "amount": 100},
],

}
*/

type EntrySetInput struct {
	Entries []LedgerEntryInput `json:"entries"`
}

// TODO: Add validation for type credit == debit
