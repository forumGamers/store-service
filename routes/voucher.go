package routes

import (
	"github.com/forumGamers/store-service/cmd"
	md "github.com/forumGamers/store-service/middlewares"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) voucherRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/voucher")

	uri.POST("/",md.Authentication,cmd.AddVoucher)

	uri.GET("/",q.GetAllVoucher)

	uri.GET("/:id",q.GetVoucherById)

	uri.DELETE("/:id",md.Authentication,cmd.DeleteVoucher)
}