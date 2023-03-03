package app

import "github.com/luizmoitinho/bookstore_users_api/logger"

func StartApplicaton() {
	app := GinRouter{}

	logger.Info("initialiazing the application")
	app.Init()

	logger.Info("mapping the api routes")
	app.MapRoutes()

	logger.Info("running into port 8081")
	app.Run(":8081")
}
