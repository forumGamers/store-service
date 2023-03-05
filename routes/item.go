package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) itemRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/item")

	uri.GET("/",c.GetAllItem)

	uri.GET("/store/:storeId",c.GetItemByStoreId)

	uri.PATCH("/change-desc/:id",c.UpdateItemDesc)

	uri.PATCH("/change-image/:id",c.UpdateItemImage)

	uri.PATCH("/add-stock/:id",c.AddStock)

	uri.PATCH("/change-price/:id",c.UpdatePrice)

	uri.PATCH("/change-name/:storeId/:id",c.UpdateName)

	uri.GET("/:slug",c.GetItemBySlug)

	uri.POST("/:storeId",c.CreateItem)
}