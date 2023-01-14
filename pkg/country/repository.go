package country

import (
	"database/sql"
)

type CountriesApex interface {
	Retrieve(countryCode string) (*Country, error)
}

type PersistentCountriesApex struct {
	db *sql.DB
}

func NewPersistentCountriesApex(db *sql.DB) *PersistentCountriesApex {
	return &PersistentCountriesApex{db: db}
}

func (p *PersistentCountriesApex) Retrieve(countryCode string) (*Country, error) {
	row := p.db.QueryRow("select sum(rangos) as cant_ips from (select *,(ips.ip_to - ips.ip_from + 1) as rangos from ips) as detailed_ips where country_code = $1", countryCode)
	country, err := NewCountryFromRow(row)
	if err == sql.ErrNoRows {
		return nil, CountryNotFound
	}
	return country, nil
}
