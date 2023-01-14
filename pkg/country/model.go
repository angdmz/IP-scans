package country

import "ipScans/pkg/utils"

type Country struct {
	IpCount uint
}

func NewCountryFromRow(row utils.Scanneable) (*Country, error) {
	country := &Country{}
	err := row.Scan(
		&country.IpCount,
	)
	if err != nil {
		return nil, err
	}
	return country, nil
}
