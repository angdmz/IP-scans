package main

import (
	"fmt"
	"ipScans/pkg/country"
	database "ipScans/pkg/db"
	"ipScans/pkg/ips"
	"ipScans/pkg/provider"
	serverouter "ipScans/pkg/router"
)

func main() {
	db, err := database.OpenDB("172.17.0.1", 5432, "ipscans", "ipscans", "ipscans")
	if err != nil {
		panic(err)
	}
	handler := serverouter.NewHandler(ips.NewIpAgenda(db), provider.NewProvidersAgenda(db), country.NewPersistentCountriesApex(db))
	router := serverouter.ServeRouter(handler)
	err = router.Run()
	if err != nil {
		fmt.Printf(err.Error())
	}
}
