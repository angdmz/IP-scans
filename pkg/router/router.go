package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	country "ipScans/pkg/country"
	"ipScans/pkg/ips"
	"ipScans/pkg/provider"
	"net/http"
)

func ServeRouter(handler Handler) *gin.Engine {
	r := gin.Default()
	handler.AddRoutes(r)
	return r
}

type Handler interface {
	AddRoutes(engine *gin.Engine)
}

type V1Handler struct {
	ipAgenda        ips.IPAgenda
	providersAgenda provider.ProvidersAgenda
	countriesApex   country.CountriesApex
}

func (h *V1Handler) AddRoutes(r *gin.Engine) {
	r.GET("/health", h.HandleHealth)
	r.GET("/api/v1/ips/:ip", h.HandleRetrieveIP)
	r.GET("/api/v1/providers", h.HandleListProviders)
	r.GET("/api/v1/countries/:country_code", h.HandleRetrieveCountry)
}

func NewHandler(ipAgenda ips.IPAgenda, providersAgenda provider.ProvidersAgenda, countriesApex country.CountriesApex) *V1Handler {
	return &V1Handler{ipAgenda: ipAgenda, providersAgenda: providersAgenda, countriesApex: countriesApex}
}

func (h *V1Handler) HandleHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *V1Handler) HandleRetrieveIP(context *gin.Context) {
	paramIp, exists := context.Params.Get("ip")
	if exists {
		ip, err := h.ipAgenda.Retrieve(paramIp)
		if err != nil {
			statusCode := http.StatusInternalServerError
			if errors.Is(err, ips.NotIPFormat) {
				statusCode = http.StatusUnprocessableEntity
			} else if errors.Is(err, ips.IpNotFound) {
				statusCode = http.StatusNotFound
			}
			context.JSON(statusCode, gin.H{"messages": [1]string{err.Error()}})
			return
		}
		context.JSON(http.StatusOK, ip)
	} else {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"messages": [1]string{"Param not found"}})
	}
}

func (h *V1Handler) HandleListProviders(context *gin.Context) {
	// sort := context.DefaultQuery("sort", "desc")
	countryCode := context.Query("country_code")
	if countryCode != "" {
		providers, err := h.providersAgenda.List(countryCode)
		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"messages": [1]string{err.Error()}})
			return
		}
		context.JSON(http.StatusOK, gin.H{"count": len(providers), "results": providers})
	} else {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"messages": [1]string{"Param not found"}})
	}
}

func (h *V1Handler) HandleRetrieveCountry(context *gin.Context) {
	countryCode, exists := context.Params.Get("country_code")
	if exists {
		var returnable any
		c, err := h.countriesApex.Retrieve(countryCode)
		statusCode := http.StatusOK
		returnable = c
		if err != nil {
			returnable = gin.H{"messages": [1]string{err.Error()}}
			if errors.Is(err, country.CountryNotFound) {
				statusCode = http.StatusNotFound
			} else {
				statusCode = http.StatusInternalServerError
			}
		}
		context.JSON(statusCode, returnable)
	} else {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"messages": [1]string{"Param not found"}})
	}
}
