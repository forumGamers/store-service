package routes

import (
	"github.com/forumGamers/store-service/cmd"
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) cartRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/cart")

	uri.POST("/:itemId",cmd.AddCart)

	uri.GET("/",c.GetCart)

	uri.DELETE("/:id",cmd.RemoveCart)
}