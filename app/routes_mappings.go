package app

import (
	"github.com/gin-gonic/gin"
	ping "github.com/luizmoitinho/bookstore_users_api/controllers/ping"
	users "github.com/luizmoitinho/bookstore_users_api/controllers/users"
)

type GinRouter struct {
	engine *gin.Engine
}

func (g *GinRouter) Init() {
	g.engine = gin.Default()
}

func (g *GinRouter) Run(address string) {
	g.engine.Run(address)
}

func (g *GinRouter) MapRoutes() {
	g.engine.GET("/ping", ping.Ping)

	g.engine.GET("/users/:user_id", users.GetUser)
	g.engine.POST("/users", users.CreateUser)
	g.engine.PUT("/users/:user_id", users.UpdateUser)
	g.engine.PATCH("/users/:user_id", users.UpdateUser)

}
