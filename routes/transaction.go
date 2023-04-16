package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) transactionRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/transaction")

	uri.POST("/:storeId/:itemId",md.Authentication,cmd.CreateTransaction)

	uri.PATCH("/:transactionId",md.Authentication,cmd.EndTransaction)

	uri.PATCH("/cancel/:transactionId",md.Authentication,cmd.CancelTransaction)

	uri.GET("/",q.GetAllTransaction)
}