package routes

import (
	"github.com/forumGamers/store-service/cmd"
	q "github.com/forumGamers/store-service/query"
	"github.com/gin-gonic/gin"
)

func (r routes) voucherRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/voucher")

	uri.POST("/",cmd.AddVoucher)

	uri.GET("/",q.GetAllVoucher)

	uri.GET("/:id",q.GetVoucherById)

	uri.DELETE("/:id",cmd.DeleteVoucher)
}