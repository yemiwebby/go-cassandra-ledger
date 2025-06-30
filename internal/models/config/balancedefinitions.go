package config

type BalanceDefinition struct {
	TimeAxis string   `yaml:"time_axis"`
	Address  []string `yaml:"addresses"`
}
