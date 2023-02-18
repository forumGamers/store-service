package routes

import (
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

	groupRoutes := r.router.Group("/api")

	r.storeRoutes(groupRoutes)

	r.router.Run(":4000")
}