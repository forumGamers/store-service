package main

import (
	c "github.com/forumGamers/store-service/config"
	r "github.com/forumGamers/store-service/routes"
)

func main(){
	c.Connection()

	r.Routes()
}