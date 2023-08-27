package main

import (
	"avito_internship/internal/config"
	"avito_internship/internal/database"
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

	//_, err = db.Create("avito_tech")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(id)

	//arr := []string{"test", "avito_tech"}
	//err = db.AddUser(arr, []string{}, 3)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = db.CreateUser()
	//if err != nil {
	//	log.Fatal(err)
	//}

	err = db.Delete("test")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to database")
}
