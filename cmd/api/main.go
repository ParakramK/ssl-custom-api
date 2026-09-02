package main

import (
	"log"
	"net/http"

	"ssl-custom-api/internal/config"
	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/sap/aging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	hanaProvider, err := hana.NewProvider(hana.Config{
		Host:     cfg.HANAHost,
		Port:     cfg.HANAPort,
		User:     cfg.HANAUser,
		Password: cfg.HANAPassword,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer hanaProvider.Close()

	log.Println("HANA provider initialized")

	repository := aging.NewRepository(hanaProvider)
	service := aging.NewService(repository)
	handler := aging.NewHandler(service)

	http.HandleFunc("/api/aging", handler.GetCustomerAging)

	log.Println("API listening on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
