package routes

import (
	"net/http"
	"os"

	md "github.com/forumGamers/store-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type routes struct {
	router *gin.Engine
}

func Routes(){

	if err := godotenv.Load() ; err != nil {
		panic(err.Error())
	}

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

	r.itemRoutes(groupRoutes)

	r.followerRoutes(groupRoutes)

	r.favoriteRoutes(groupRoutes)

	r.transactionRoutes(groupRoutes)

	r.cartRoutes(groupRoutes)

	r.voucherRoutes(groupRoutes)

	port := os.Getenv("PORT")

	if port == "" {
		port = "4000"
	}

	r.router.Run(":"+port)
}