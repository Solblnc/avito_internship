package main

import (
	"avito_internship/internal/config"
	"avito_internship/internal/database"
	"avito_internship/internal/server"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDataBase(database.Config{
		Host:     cfg.Config.Host,
		Port:     cfg.Config.Port,
		User:     cfg.Config.User,
		DBName:   cfg.Config.DBName,
		Password: cfg.Config.Password,
		SSLMode:  cfg.Config.SSLMode,
	})
	if err != nil {
		log.Fatal(err)
	}

	newServer, err := server.NewServer(cfg.Port.Port, db)
	if err != nil {
		log.Fatal(err)
	}
	newServer.SetupHandlers()
	newServer.Run()

	log.Println("Connected to database")
}
