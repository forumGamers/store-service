package routes

import (
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) followerRoutes(rg *gin.RouterGroup){

	uri := rg.Group("/follower")

	uri.GET("/:storeId",q.GetStoreFollower)
}