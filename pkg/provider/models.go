package provider

import "ipScans/pkg/utils"

type Provider struct {
	Name     string `json:"name"`
	IpAmount uint   `json:"ip_amount"`
}

func NewProviderFromRow(row utils.Scanneable) (*Provider, error) {
	provider := &Provider{}
	err := row.Scan(
		&provider.Name,
		&provider.IpAmount,
	)
	if err != nil {
		return nil, err
	}
	return provider, nil
}
