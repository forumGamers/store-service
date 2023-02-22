package middlewares

import "github.com/gin-gonic/gin"

func AuthorizeAdmin(c *gin.Context) {
	role := c.Request.Header.Get("role")

	if role != "Admin" {
		panic("Forbidden")
	}

	c.Next()
}