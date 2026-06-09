package router

import (
	// "fmt"
	// "neneloga/internal/handlers"

	"neneloga/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/",handlers.Home)

	// health
	r.GET("/health", handlers.Health)

	r.GET("/chuck",handlers.ChuckNorris)

	// r.GET("/ping", handlers.Ping)

	// user := r.Group("/users")
	// {
	//     user.GET("/", handler.GetUsers)
	//     user.POST("/", handler.CreateUser)
	// }

	r.NoRoute(handlers.NotFound)

	return r
}
