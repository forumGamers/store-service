package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) store_status_routes(rg *gin.RouterGroup ){

	uri := rg.Group("/store_status")

	uri.POST("/",md.AuthorizeAdmin,cmd.CreateStoreStatus)

	uri.GET("/",q.GetAllStoreStatus)

	uri.PATCH("/change-name/:id",md.AuthorizeAdmin,cmd.UpdateStoreStatusName)

	uri.PATCH("/change-exp/:id",md.AuthorizeAdmin,cmd.UpdateStoreStatusExp)
}