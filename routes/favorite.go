package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) favoriteRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/favorite")

	uri.Use(md.Authentication)

	uri.GET("/",q.GetMyFavorite)

	uri.POST("/:itemId",cmd.AddFavorite)

	uri.DELETE("/:id",cmd.RemoveFavorite)
}