package routes

import (
	"github.com/forumGamers/store-service/cmd"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) cartRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/cart")

	uri.POST("/:itemId",cmd.AddCart)

	uri.GET("/",q.GetCart)

	uri.DELETE("/:id",cmd.RemoveCart)
}