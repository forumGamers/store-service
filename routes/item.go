package routes

import (
	"github.com/forumGamers/store-service/cmd"
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) itemRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/item")

	uri.GET("/",c.GetAllItem)

	uri.GET("/store/:storeId",c.GetItemByStoreId)

	uri.PATCH("/change-desc/:id",cmd.UpdateItemDesc)

	uri.PATCH("/change-image/:id",cmd.UpdateItemImage)

	uri.PATCH("/add-stock/:id",cmd.AddStock)

	uri.PATCH("/change-price/:id",cmd.UpdatePrice)

	uri.PATCH("/change-name/:id",cmd.UpdateName)

	uri.PATCH("/change-discount/:id",cmd.UpdateItemDiscount)

	uri.GET("/:slug",c.GetItemBySlug)

	uri.POST("/:storeId",cmd.CreateItem)
}