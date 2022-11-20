package main

import (
	"github.com/luizmoitinho/bookstore_users_api/app"
	"github.com/luizmoitinho/bookstore_users_api/config"
)

func main() {
	config.Load(".env")

	app.StartApplicaton()
}
