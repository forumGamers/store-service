package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/store")

	uri.POST("/",md.Authentication,cmd.CreateStore)

	uri.GET("/",q.GetAllStores)

	uri.GET("/name",q.GetStoreName)

	uri.GET("/my-store",md.Authentication,q.GetMyStore)

	uri.PATCH("/change-name",md.Authentication,cmd.UpdateStoreName)

	uri.PATCH("/change-desc",md.Authentication,cmd.UpdateStoreDesc)

	uri.PATCH("/change-image",md.Authentication,cmd.UpdateStoreImage)

	uri.PATCH("/change-background",md.Authentication,cmd.UpdateStoreBg)

	uri.PATCH("/deactived",md.Authentication,cmd.DeactiveStore)

	uri.PATCH("/reactived",md.Authentication,cmd.ReactivedStore)

	uri.GET("/:id",q.GetStoreById)
}