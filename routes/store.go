package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/store")

	uri.POST("/",c.CreateStore)
}