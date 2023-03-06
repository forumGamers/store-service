package routes

import (
	c "github.com/forumGamers/store-service/controllers"
	"github.com/gin-gonic/gin"
)

func (r routes) voucherRoutes(rg *gin.RouterGroup){
	uri := rg.Group("/voucher")

	uri.POST("/",c.AddVoucher)
}