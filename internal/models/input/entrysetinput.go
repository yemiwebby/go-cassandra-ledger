package input

type LedgerAddress struct {
	LegalEntity string `json:"legal_entity"`
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Currency    string `json:"currency"`
	AccountID   string `json:"account_id"`
}

type LedgerEntryInput struct {
	Address       LedgerAddress `json:"address"`
	Type          string        `json:"type"` // credit/debit
	Amount        float64       `json:"amount"`
	Description   string        `json:"description"`
	Timestamp     int64         `json:"timestamp"`
	ReportingTime *int64        `json:"reporting_ts,omitempty"`
}

type EntrySetInput struct {
	Entries []LedgerEntryInput `json:"entries"`
}
