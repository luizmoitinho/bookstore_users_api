package users_db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/luizmoitinho/bookstore_users_api/config"
	"github.com/luizmoitinho/bookstore_users_api/logger"
)

const (
	DATASOURCE_NAME = "%s:%s@tcp(%s:%s)/%s?charset=%s"
)

type MySQL struct {
	Client *sql.DB
}

func Connect() MySQL {
	var err error
	var _mysql = MySQL{}

	_mysql.Client, err = sql.Open("mysql", fmt.Sprintf(
		DATASOURCE_NAME,
		config.Propertie.DATABASE.USER,
		config.Propertie.DATABASE.PASS,
		config.Propertie.DATABASE.HOST,
		config.Propertie.DATABASE.PORT,
		config.Propertie.DATABASE.NAME,
		config.Propertie.DATABASE.COLLECTION,
	))
	if err != nil {
		logger.Error("error when trying connect to database", err)
		panic(err)
	}

	if err = _mysql.Client.Ping(); err != nil {
		panic(err)
	}

	logger.Info("database succesfully connected connected")
	return _mysql
}
