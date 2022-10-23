package app

import "github.com/luizmoitinho/bookstore_users_api/controllers"

func mapRoutes() {
	router.GET("/ping", controllers.Ping)
}
