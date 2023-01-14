package provider

import (
	"database/sql"
)

type PersistentProvidersAgenda struct {
	db *sql.DB
}

func NewProvidersAgenda(db *sql.DB) *PersistentProvidersAgenda {
	return &PersistentProvidersAgenda{db: db}
}

func (agenda *PersistentProvidersAgenda) List(countryCode string) ([]*Provider, error) {
	rows, err := agenda.db.Query("select isp, sum(rangos) as cant_ips from (select *,(ips.ip_to - ips.ip_from + 1) as rangos from ips) as detailed_ips where country_code = $1 group by isp order by cant_ips desc limit 10;", countryCode)
	if err != nil {
		return nil, err
	}
	var providers []*Provider
	for rows.Next() {
		prov, err := NewProviderFromRow(rows)
		if err != nil {
			return nil, err
		}
		providers = append(providers, prov)
	}
	return providers, nil
}

type ProvidersAgenda interface {
	List(countryCode string) ([]*Provider, error)
}
