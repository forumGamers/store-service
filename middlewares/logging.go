package middlewares

import (
	"fmt"
	"time"

	h "github.com/forumGamers/store-service/helper"
	"github.com/forumGamers/store-service/loaders"
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
)

func Logging(c *gin.Context){
	defer func() {
		id := h.GetUser(c).Id

		responseTime := c.MustGet("start").(time.Time)

		if err := loaders.GetDb().Model(m.Log{}).Create(&m.Log{
			Path: c.Request.URL.Path,
			UserId: id ,
			Method: c.Request.Method,
			StatusCode: c.Writer.Status(),
			Origin: c.Request.Header.Get("Origin"),
			ResponseTime: int(time.Since(responseTime).Milliseconds()),
		}).Error ; err != nil {
			fmt.Println(err)
			return
		}
	}()
	c.Next()
}