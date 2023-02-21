package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) itemRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/item")

	uri.POST("/",c.CreateItem)
}