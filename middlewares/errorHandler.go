package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	defer func(){
		msg :=  recover()
		s := http.StatusInternalServerError
		if msg == "Data not found"{
			s = http.StatusNotFound
		}else if msg == "Invalid data" {
			s = http.StatusBadRequest
		}else if msg == "Forbidden" {
			s = http.StatusForbidden
		}else {
			msg = "Internal Server Error"
		}
		c.AbortWithStatusJSON(s,gin.H{"message":msg})
	}()
	c.Next()
}