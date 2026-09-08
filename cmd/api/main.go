package main

import (
	"log"

	"ssl-custom-api/internal/app"
	"ssl-custom-api/internal/app/auth"
	"ssl-custom-api/internal/config"
	"ssl-custom-api/internal/providers/hana"
	"ssl-custom-api/internal/providers/mysql"
	"ssl-custom-api/internal/providers/postgres"
	"ssl-custom-api/internal/router"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	jwtService, err_jwt := auth.NewJWT(cfg.JWTSecret)
	if err_jwt != nil {
		log.Fatal(err_jwt)
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

	mysqlProvider, err := mysql.NewProvider(mysql.Config{
		Host:     cfg.MySQLHost,
		Port:     cfg.MySQLPort,
		User:     cfg.MySQLUser,
		Password: cfg.MySQLPassword,
		Database: cfg.MySQLDatabase,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer mysqlProvider.Close()

	postgresProvider, err := postgres.NewProvider(postgres.Config{
		Host:     cfg.PGHost,
		Port:     cfg.PGPort,
		User:     cfg.PGUser,
		Password: cfg.PGPassword,
		Database: cfg.PGDatabase,
		SSLMode:  cfg.PGSSLMode,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer postgresProvider.Close()

	log.Println("HANA provider initialized")
	log.Println("MySQL provider initialized")
	log.Println("PostgreSQL provider initialized")

	handlers := app.New(hanaProvider, mysqlProvider, postgresProvider, jwtService, cfg.APIKey)
	r := router.Setup(handlers)

	log.Println("API listening on :8080")

	if err := r.Listen(":8080"); err != nil {
		log.Fatal(err)
	}
}
