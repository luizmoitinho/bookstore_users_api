package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Propertie Properties

type Properties struct {
	ENV         string
	CONTAINER   bool
	PORT        int
	LOCATION    string
	DATE_LAYOUT string
	DATABASE    DataBase
}

type DataBase struct {
	HOST       string
	PORT       string
	NAME       string
	USER       string
	PASS       string
	COLLECTION string
}

func Load(path string) {
	var err error

	if err = godotenv.Load(path); err != nil {
		log.Fatal("config.Load(): ", err)
	}

	Propertie.DATE_LAYOUT = os.Getenv("DATE_LAYOUT")
	Propertie.LOCATION = os.Getenv("LOCATION")
	Propertie.ENV = os.Getenv("ENV")
	if Propertie.PORT, err = strconv.Atoi(os.Getenv("PORT")); err != nil {
		Propertie.PORT = 8080
		log.Printf("cannot convert PORT at .env: %v", err)
	}

	//mysql
	Propertie.DATABASE.HOST = os.Getenv("DB_HOST")
	Propertie.DATABASE.PORT = os.Getenv("DB_PORT")
	Propertie.DATABASE.NAME = os.Getenv("DB_NAME")
	Propertie.DATABASE.USER = os.Getenv("DB_USER")
	Propertie.DATABASE.PASS = os.Getenv("DB_PASS")
	Propertie.DATABASE.COLLECTION = os.Getenv("DB_COLLECTION")

}
