package helper

import (
	m "github.com/forumGamers/store-service/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GetUser(c *gin.Context) m.User {
	claim := c.MustGet("user").(jwt.MapClaims)
	var user m.User

	for key, val := range claim {
		switch key {
		case "email":
			user.Email = val.(string)
		case "fullName":
			user.Fullname = val.(string)
		case "iat":
			user.Iat = int(val.(float64))
		case "id":
			user.Id = int(val.(float64))
		case "isVerified":
			user.IsVerified = val.(bool)
		case "phoneNumber":
			user.PhoneNumber = val.(string)
		case "username":
			user.Username = val.(string)
		case "StoreId" :
			user.StoreId = int(val.(float64))
		case "role" :
			user.Role = val.(string)
		case "point" :
			user.Point = int(val.(float64))
		case "exp" :
			user.Exp = int(val.(float64))
		}
	}
	return user
}