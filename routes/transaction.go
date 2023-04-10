package routes

import (
	"github.com/forumGamers/store-service/cmd"
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) transactionRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/transaction")

	uri.POST("/:storeId/:itemId",cmd.CreateTransaction)

	uri.PATCH("/:transactionId",cmd.EndTransaction)

	uri.PATCH("/cancel/:transactionId",cmd.CancelTransaction)

	uri.GET("/",c.GetAllTransaction)
}