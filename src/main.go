package main

import (
	_ "MeetYou/src/Game_Logic"
	internet "MeetYou/src/Internet_interface"
	mysqlConn "MeetYou/src/MySQL_Connection"
)

func main(){
	internet.InitLog()
	mysqlConn.Connect()
	internet.Start()
}
