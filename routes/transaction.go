package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) transactionRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/transaction")

	uri.POST("/:storeId/:itemId",c.CreateTransaction)

	uri.PATCH("/:transactionId",c.EndTransaction)

	uri.PATCH("/cancel/:transactionId",c.CancelTransaction)
}