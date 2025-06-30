package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func LoadResolvedBalanceDefinitions(
	addressConfigPath, balanceDefinitionPath string,
) (map[string]ResolvedBalanceDefinition, error) {

	addressBytes, err := os.ReadFile(addressConfigPath)
	if err != nil {
		return nil, fmt.Errorf("reading address config: %w", err)
	}

	addressMap := map[string]LedgerAddress{}
	if err := yaml.Unmarshal(addressBytes, &addressMap); err != nil {
		return nil, fmt.Errorf("unmarshalling address config: %w", err)
	}

	balanceBytes, err := os.ReadFile(balanceDefinitionPath)
	if err != nil {
		return nil, fmt.Errorf("reading balance definitions: %w", err)
	}

	balanceMap := map[string]BalanceDefinition{}
	if err := yaml.Unmarshal(balanceBytes, &balanceMap); err != nil {
		return nil, fmt.Errorf("unmarshalling balance definitions: %w", err)
	}

	resolved := map[string]ResolvedBalanceDefinition{}
	for name, def := range balanceMap {
		var resolvedAddresses []LedgerAddress
		for _, addrKey := range def.Address {
			addr, ok := addressMap[addrKey]
			if !ok {
				return nil, fmt.Errorf("address reference %q not found", addrKey)
			}
			resolvedAddresses = append(resolvedAddresses, addr)
		}
		resolved[name] = ResolvedBalanceDefinition{
			Name:      name,
			TimeAxis:  def.TimeAxis,
			Addresses: resolvedAddresses,
		}
	}

	return resolved, nil
}
