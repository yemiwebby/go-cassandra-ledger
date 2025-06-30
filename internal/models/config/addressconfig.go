package config

type LedgerAddress struct {
	LegalEntity string `yaml:"legal_entity"`
	Namespace   string `yaml:"namespace"`
	Name        string `yaml:"name"`
	Currency    string `yaml:"currency"`
	AccountID   string `yaml:"account_id"`
}
