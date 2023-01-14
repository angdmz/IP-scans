package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"ipScans/pkg/country"
	database "ipScans/pkg/db"
	"ipScans/pkg/ips"
	"ipScans/pkg/provider"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Health struct {
	Status string `json:"status"`
}

func (h *Health) isOk() bool {
	return h.Status == "ok"
}

type ErrorMessage struct {
	Messages []string `json:"messages"`
}

type RetrieveIp struct {
	ips.Ip
}

type RetrieveProvider struct {
	provider.Provider
}

type ListProvider struct {
	Count   int                `json:"count"`
	Results []RetrieveProvider `json:"results"`
}

type RetrieveCountry struct {
	country.Country
}

func AssertRetrieveIpFromBytes(t *testing.T, bytes []byte) {
	retrievedIp := &RetrieveIp{}
	err := json.Unmarshal(bytes, &retrievedIp)
	assert.Nil(t, err)
	assert.Equal(t, retrievedIp.ProxyType, "PUB")
	assert.Equal(t, retrievedIp.CountryCode, "TH")
	assert.Equal(t, retrievedIp.CountryName, "Thailand")
	assert.Equal(t, retrievedIp.RegionName, "Chiang Rai")
	assert.Equal(t, retrievedIp.CityName, "Pa Daet")
	assert.Equal(t, retrievedIp.Isp, "TOT Public Company Limited")
	assert.Equal(t, retrievedIp.Domain, "tot.co.th")
	assert.Equal(t, retrievedIp.UsageType, "ISP/MOB")
	assert.Equal(t, retrievedIp.Asn, "23969")
	assert.Equal(t, retrievedIp.As, "TOT Public Company Limited")
}

func generateRouter() *gin.Engine {
	db, err := database.OpenDB("172.17.0.1", 5432, "ipscans", "ipscans", "ipscans")
	if err != nil {
		panic(err)
	}
	handler := NewHandler(ips.NewIpAgenda(db), provider.NewProvidersAgenda(db), country.NewPersistentCountriesApex(db))
	router := ServeRouter(handler)
	return router
}

func TestHealthRoute(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var pong Health
	err := json.Unmarshal(w.Body.Bytes(), &pong)
	assert.Nil(t, err)
	assert.True(t, pong.isOk())
}

func TestRetrieveIpRoute(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ips/1.0.128.137", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	AssertRetrieveIpFromBytes(t, w.Body.Bytes())
}

func TestRetrieveIpNotExistingRoute(t *testing.T) {
	router := generateRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ips/0.0.0.0", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	var errors ErrorMessage
	err := json.Unmarshal(w.Body.Bytes(), &errors)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(errors.Messages))
	assert.Equal(t, "ip not found", errors.Messages[0])
}

func TestRetrieveBullshitIpRoute(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ips/sarasa", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	var errors ErrorMessage
	err := json.Unmarshal(w.Body.Bytes(), &errors)
	assert.Nil(t, err)
	assert.Equal(t, len(errors.Messages), 1)
	assert.Equal(t, errors.Messages[0], "incorrect IP format")
}

func TestListProviders(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/providers?country_code=CH&sort=desc", nil)
	router.ServeHTTP(w, req)
	var list ListProvider
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &list)
	assert.Nil(t, err)
	assert.Equal(t, len(list.Results), list.Count)
	assert.Equal(t, list.Count, 10)
	assert.Equal(t, list.Results[0].Name, "Microsoft Corporation")
	assert.Equal(t, list.Results[1].Name, "Sunrise GmbH")
	assert.Equal(t, list.Results[2].Name, "RapidSeedbox Ltd")
	assert.Equal(t, list.Results[3].Name, "Rosite Equipment SRL")
	assert.Equal(t, list.Results[4].Name, "Google LLC")
	assert.Equal(t, list.Results[5].Name, "Cloud Innovation Ltd")
	assert.Equal(t, list.Results[6].Name, "Swisscom AG")
	assert.Equal(t, list.Results[7].Name, "Valaiscom AG")
	assert.Equal(t, list.Results[8].Name, "Private Layer Inc")
	assert.Equal(t, list.Results[9].Name, "Bluewin is an LIR and ISP in Switzerland.")
	assert.Equal(t, list.Results[0].IpAmount, uint(101))
	assert.Equal(t, list.Results[1].IpAmount, uint(99))
	assert.Equal(t, list.Results[2].IpAmount, uint(93))
	assert.Equal(t, list.Results[3].IpAmount, uint(84))
	assert.Equal(t, list.Results[4].IpAmount, uint(72))
	assert.Equal(t, list.Results[5].IpAmount, uint(62))
	assert.Equal(t, list.Results[6].IpAmount, uint(57))
	assert.Equal(t, list.Results[7].IpAmount, uint(52))
	assert.Equal(t, list.Results[8].IpAmount, uint(45))
	assert.Equal(t, list.Results[9].IpAmount, uint(30))
}

func TestListNonExistingCountryProviders(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/providers?country_code=poronga&sort=desc", nil)
	router.ServeHTTP(w, req)
	var list ListProvider
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &list)
	assert.Nil(t, err)
	assert.Equal(t, len(list.Results), list.Count)
	assert.Equal(t, len(list.Results), 0)
}

func TestListCountryIps(t *testing.T) {
	router := generateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/countries/CH", nil)
	router.ServeHTTP(w, req)
	var retrieveCountry RetrieveCountry
	assert.Equal(t, http.StatusOK, w.Code)
	err := json.Unmarshal(w.Body.Bytes(), &retrieveCountry)
	assert.Nil(t, err)
	assert.Equal(t, retrieveCountry.IpCount, uint(1056))
}
