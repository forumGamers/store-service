package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func(){
		msg :=  recover()
		s := http.StatusInternalServerError
		if msg == nil {
			return
		}
		switch msg {
		case "record not found" :
			msg = "Data not found"
		case "Data not found" :
			s = http.StatusNotFound
			break
		case "Failed to parse image" :
			s = http.StatusBadRequest
			break
		case "Forbidden":
			s = http.StatusForbidden
			break
		case "Invalid params" :
			s = http.StatusBadRequest
			break
		case "You have not store yet" : 
			s = http.StatusBadRequest
			break
		case "Invalid data":
			s = http.StatusBadRequest
			break
		case "Stock is not enough" :
			s = http.StatusBadRequest
			break
		case "voucher is not registered" :
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
			break
		case "Unauthorize" :
			s = http.StatusUnauthorized
			break
		case "Conflict":
			s = http.StatusConflict
			break
		case "file cannot be larger than 10 mb":
			s = http.StatusBadRequest
			break
		case "file type is not supported":
			s = http.StatusBadRequest
			break
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