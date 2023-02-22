package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	md "github.com/forumGamers/store-service/middlewares"
	"github.com/gin-gonic/gin"
)

func (r routes) store_status_routes(rg *gin.RouterGroup ){

	uri := rg.Group("/store_status")

	uri.POST("/",md.AuthorizeAdmin,c.CreateStoreStatus)

	uri.GET("/",c.GetAllStoreStatus)

	uri.PATCH("/change-name/:id",md.AuthorizeAdmin,c.UpdateStoreStatusName)
}