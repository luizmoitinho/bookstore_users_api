package users_db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luizmoitinho/bookstore_users_api/config"
)

const (
	DATASOURCE_NAME = "%s:%s@tcp(%s:%s)/%s?charset=%s"
)

type MySQL struct {
	DB *sql.DB
}

func Connect() MySQL {
	var err error
	var _mysql = MySQL{}

	_mysql.DB, err = sql.Open("mysql", fmt.Sprintf(
		DATASOURCE_NAME,
		config.Propertie.DATABASE.USER,
		config.Propertie.DATABASE.PASS,
		config.Propertie.DATABASE.HOST,
		config.Propertie.DATABASE.PORT,
		config.Propertie.DATABASE.NAME,
		config.Propertie.DATABASE.COLLECTION,
	))
	if err != nil {
		panic(err)
	}

	if err = _mysql.DB.Ping(); err != nil {
		panic(err)
	}

	log.Println("database succesfully connected connected")
	return _mysql
}
