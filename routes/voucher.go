package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) voucherRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/voucher")

	uri.POST("/",c.AddVoucher)

	uri.GET("/",c.GetAllVoucher)

	uri.GET("/:id",c.GetVoucherById)

	uri.DELETE("/:id",c.DeleteVoucher)
}