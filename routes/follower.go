package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) followerRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/follower")

	uri.GET("/:storeId",q.GetStoreFollower)

	uri.POST("/:storeId",md.Authentication,cmd.FollowStoreById)
}