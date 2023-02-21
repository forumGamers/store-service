package routes

import (
	"net/http"

	md "github.com/forumGamers/store-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
)

type routes struct {
	router *gin.Engine
}

func Routes(){
	r := routes { router: gin.Default() }

	r.router.Use(cors.Default())
	
	r.router.Use(logger.SetLogger())

	r.router.Use(md.ErrorHandler)

	//testing connection
	r.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK,gin.H{"message":"pong"})
	})

	groupRoutes := r.router.Group("/api")

	r.storeRoutes(groupRoutes)

	r.store_status_routes(groupRoutes)

	r.router.Run(":4000")
}