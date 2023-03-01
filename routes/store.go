package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/store")

	uri.POST("/",c.CreateStore)

	uri.GET("/",c.GetAllStores)

	uri.PATCH("/change-name",c.UpdateStoreName)

	uri.PATCH("/change-desc",c.UpdateStoreDesc)

	uri.PATCH("/change-image",c.UpdateStoreImage)

	uri.PATCH("/deactived",c.DeactiveStore)

	uri.PATCH("/reactived",c.ReactivedStore)

	uri.GET("/:id",c.GetStoreById)
}