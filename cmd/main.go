package main

import (
	"avito_internship/internal/config"
	"avito_internship/internal/database"
	"fmt"
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

	//_, err = db.Create("another")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(id)

	//arrDelete := []string{"test", "avito_tech"}
	//arrAdd := []string{"test", "avito_tech", "another"}

	//err = db.AddUser(arr, []string{}, 10)
	//err = db.AddUser(arrAdd, arrDelete, 5)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println(db.GetActualSegments(1))

	//err = db.CreateUser()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//err = db.Delete("test")
	//if err != nil {
	//	log.Fatal(err)
	//}

	log.Println("Connected to database")
}
