package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) cartRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/cart")

	uri.POST("/:itemId",c.AddCart)

	uri.GET("/",c.GetCart)

	uri.DELETE("/:id",c.RemoveCart)
}