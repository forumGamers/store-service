package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) store_status_routes(rg *gin.RouterGroup ){

	uri := rg.Group("/store_status")

	uri.POST("/",c.CreateStoreStatus)

	uri.GET("/",c.GetAllStoreStatus)
}