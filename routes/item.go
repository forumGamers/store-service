package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) itemRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/item")

	uri.POST("/",c.CreateItem)

	uri.GET("/",c.GetAllItem)

	uri.GET("/store/:storeId",c.GetItemByStoreId)

	uri.PATCH("/change-desc/:id",c.UpdateItemDesc)

	uri.PATCH("/change-name/:id",c.UpdateItemName)

	uri.PATCH("/change-image/:id",c.UpdateItemImage)

	uri.PATCH("/add-stock/:id",c.AddStock)

	uri.GET("/:slug",c.GetItemBySlug)
}