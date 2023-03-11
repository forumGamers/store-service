package main

import (
	c "github.com/forumGamers/store-service/config"
	r "github.com/forumGamers/store-service/routes"
)

func main(){
	c.Connection()

	r.Routes()

	//buat cron untuk ngurangi tanggal voucher

	//buat cron untuk otomatis batalin transaksi yang udh lbh 24jam
}