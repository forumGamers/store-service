package main

import (
	c "github.com/forumGamers/store-service/config"
	r "github.com/forumGamers/store-service/routes"
	s "github.com/forumGamers/store-service/services"
	"github.com/robfig/cron/v3"
)

func main(){
	c.Connection()

	r.Routes()

	c := cron.New()

	c.AddFunc("0 0 * * *",s.DecrementVoucherPeriod)

	c.Start()

	select{}

	//buat cron untuk otomatis batalin transaksi yang udh lbh 24jam
}