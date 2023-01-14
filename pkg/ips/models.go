package ips

import (
	"github.com/google/uuid"
	"ipScans/pkg/utils"
	"net"
	"strconv"
	"strings"
)

type Ip struct {
	Addr        net.IP `json:"addr"`
	ProxyType   string `json:"proxy_type"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	RegionName  string `json:"region_name"`
	CityName    string `json:"city_name"`
	Isp         string `json:"isp"`
	Domain      string `json:"domain"`
	UsageType   string `json:"usage_type"`
	Asn         string `json:"asn"`
	As          string `json:"as"`
	id          uuid.UUID
}

func NewIpFromSqlRow(requestIp net.IP, row utils.Scanneable) (*Ip, error) {
	ip := &Ip{}
	err := row.Scan(
		&ip.ProxyType,
		&ip.CountryCode,
		&ip.CountryName,
		&ip.RegionName,
		&ip.CityName,
		&ip.Isp,
		&ip.Domain,
		&ip.UsageType,
		&ip.Asn,
		&ip.As,
	)
	ip.Addr = requestIp
	if err != nil {
		return nil, err
	}
	return ip, nil
}

func ip2int(ip net.IP) uint32 {
	bits := strings.Split(ip.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	// left shifting 24,16,8,0 and bitwise OR

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return uint32(sum)
}
