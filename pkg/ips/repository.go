package ips

import (
	"database/sql"
	"net"
)

type PersistentIpAgenda struct {
	db *sql.DB
}

func NewIpAgenda(db *sql.DB) *PersistentIpAgenda {
	return &PersistentIpAgenda{db: db}
}

func (agenda *PersistentIpAgenda) Retrieve(paramIp string) (*Ip, error) {
	netIp := net.ParseIP(paramIp)
	if netIp == nil {
		return nil, NotIPFormat
	}
	decimalIp := ip2int(netIp)
	row := agenda.db.QueryRow("SELECT proxy_type, country_code, country_name, region_name, city_name, isp, domain, usage_type, asn, ips.as FROM ips WHERE $1 between ip_from and ip_to", decimalIp)
	ip, err := NewIpFromSqlRow(netIp, row)
	if err == sql.ErrNoRows {
		return nil, IpNotFound
	} else if err != nil {
		return nil, err
	}
	return ip, nil
}

type IPAgenda interface {
	Retrieve(paramIp string) (*Ip, error)
}
