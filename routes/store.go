package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/store")

	uri.POST("/",c.CreateStore)

	uri.GET("/",c.GetAllStores)

	uri.PATCH("/change-name/:id",c.UpdateStoreName)

	uri.PATCH("/change-desc/:id",c.UpdateStoreDesc)

	uri.GET("/:id",c.GetStoreById)
}