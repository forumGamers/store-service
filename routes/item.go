package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) itemRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/item")

	uri.GET("/",q.GetAllItem)

	uri.POST("/",md.Authentication,cmd.CreateItem)

	uri.GET("/store/:storeId",q.GetItemByStoreId)

	uri.GET("/list-slug",q.GetListSlug)

	uri.GET("/list-slug/:storeId",q.GetItemSlugByStoreId)

	uri.PATCH("/change-desc/:id",md.Authentication,cmd.UpdateItemDesc)

	uri.PATCH("/change-image/:id",md.Authentication,cmd.UpdateItemImage)

	uri.PATCH("/add-stock/:id",md.Authentication,cmd.AddStock)

	uri.PATCH("/change-price/:id",md.Authentication,cmd.UpdatePrice)

	uri.PATCH("/change-name/:id",md.Authentication,cmd.UpdateName)

	uri.PATCH("/change-discount/:id",md.Authentication,cmd.UpdateItemDiscount)

	uri.GET("/:slug",q.GetItemBySlug)
}