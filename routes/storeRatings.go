package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	"github.com/gin-gonic/gin"
)

func (r routes) storeRatingsRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/store-rate")

	uri.POST("/:storeId",md.Authentication,cmd.RateStore)
}