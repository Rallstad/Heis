package main

import (
	//. "./elev"
	. "./driver"
	. "./orders"
	. "./statemachine"
	//. "./network"
	//. "fmt"
	. "time"
)

func main() {
	Elev_init()
	Init_orders()
	Sleep(2 * Second)

	SM()

}
