package main

import (
	//. "./elev"
	. "./driver"
	. "./orders"
	//. "./network"
	//. "fmt"
	//. "time"
)

func main() {
	Elev_init()
	for {
		
		//Test()
		Register_order()
		Set_light()
		Print_ext_orders()
		
		/*Set_motor_dir(1)


		Elev_set_floor_indicator(Elev_get_floor_sensor_signal())*/

	}
}
