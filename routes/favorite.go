package routes

import (
	"github.com/forumGamers/store-service/cmd"
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) favoriteRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/favorite")

	uri.GET("/",c.GetMyFavorite)

	uri.POST("/:itemId",cmd.AddFavorite)

	uri.DELETE("/:id",cmd.RemoveFavorite)
}