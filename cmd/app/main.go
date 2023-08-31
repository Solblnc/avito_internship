package main

import (
	_ "avito_internship/docs"
	"avito_internship/internal/config"
	"avito_internship/internal/database"
	"avito_internship/internal/server"
	"avito_internship/internal/service"
	"log"
)

// @title Dynamic User Segmentation API
// @version         1.0
// @description     User segmentation service for avito

// @contact.name   Emil Shayhulov

// @host      localhost:8080
// @BasePath  /api/v1

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

	Service := service.NewService(db)

	httpHandler := server.NewServer(Service)
	if err = httpHandler.Serve(); err != nil {
		log.Fatal(err)
	}

}
