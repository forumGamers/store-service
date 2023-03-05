package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) followerRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/follower")

	uri.GET("/:storeId",c.GetStoreFollower)
}