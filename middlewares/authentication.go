package middlewares

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

func getSecretKey() (string,error){
	if err := godotenv.Load() ; err != nil {
		return "",errors.New("Failed to load env")
	}

	return os.Getenv("SECRET"),nil
}

func Authentication(c *gin.Context){
	access_token := c.Request.Header.Get("access_token")

	if access_token == "" {
		panic("Invalid token")
	}

	secret,err := getSecretKey()

	if err != nil {
		panic(err.Error())
	}

	claim := jwt.MapClaims{}

	token,err := jwt.ParseWithClaims(access_token,claim,func(t *jwt.Token)(interface{},error){
		return []byte(secret),nil
	})

	if err != nil {
		panic("signature is invalid")
	}

	if !token.Valid {
		panic("Invalid token")
	}

	c.Set("user",claim)

	c.Next()
}