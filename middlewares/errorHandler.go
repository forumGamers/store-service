package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func ErrorHandler(c *gin.Context) {
	defer func(){
		msg :=  recover()
		s := http.StatusInternalServerError
		if msg == nil {
			return
		}
		switch msg {
		case gorm.ErrRecordNotFound :
			msg = "Data not found"
		case "Data not found" :
			s = http.StatusNotFound
			break
		case "Forbidden":
			s = http.StatusForbidden
			break
		case "Invalid data":
			s = http.StatusBadRequest
			break
		case "name is already use" :
			s = http.StatusConflict
			break
		case "you already have a store" :
			s = http.StatusConflict
			break
		case "Bad Gateway" :
			s = http.StatusBadGateway
		default :
			fmt.Println(msg)
			msg = "Internal Server Error"
			break
		}
		c.AbortWithStatusJSON(s,gin.H{"message":msg})
		return
	}()
	c.Next()
}