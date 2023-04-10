package routes

import (
	"github.com/forumGamers/store-service/cmd"
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/store")

	uri.POST("/",cmd.CreateStore)

	uri.GET("/",c.GetAllStores)

	uri.GET("/name",c.GetStoreName)

	uri.PATCH("/change-name",cmd.UpdateStoreName)

	uri.PATCH("/change-desc",cmd.UpdateStoreDesc)

	uri.PATCH("/change-image",cmd.UpdateStoreImage)

	uri.PATCH("/change-background",cmd.UpdateStoreBg)

	uri.PATCH("/deactived",cmd.DeactiveStore)

	uri.PATCH("/reactived",cmd.ReactivedStore)

	uri.GET("/:id",c.GetStoreById)
}