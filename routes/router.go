package routes

import "github.com/gin-gonic/gin"

type routes struct {
	router *gin.Engine
}

// func Routes(){
// 	r := routes { router: gin.Default()}

// 	groupRoutes := r.router.RouterGroup("/api")
// }