package statemachine

import (
	. "../driver"
	"../orders"
	. "fmt"
	. "time"
)

type Elevator struct {
	Floor int
	Dir   Elev_dir
	Order [10]int
}

var Elev Elevator

func SM() {

	order_channel := make(chan int)
	position_channel := make(chan int)

	go Elevator_position(position_channel)
	go Check_order(order_channel)
	go orders.Register_order()

	Event_manager(order_channel, position_channel)
	/*for {
	select {
	case floor := <-position_channel:
		//Println("Current floor: ", floor)

	}*/
	//}

}

func Elevator_position(position_channel chan int) {
	floor := Elev_get_floor_sensor_signal()
	if floor != -1 {
		position_channel <- floor
		Elev_set_floor_indicator(floor)
	}

}

func Should_stop(order_channel chan int, floor int, dir Elev_dir) int {
	//rintln("jeg har liten tiss")
	if dir == 1 {
		for i := 0; i < 3; i += 2 {
			if orders.ORDER_MATRIX[floor][i] == 1 {
				Println("liten tiss")
				orders.ORDER_MATRIX[floor][i] = 0
				return 1
			}
		}
	} else if dir == -1 {
		for i := 1; i < 3; i++ {
			if orders.ORDER_MATRIX[floor][i] == 1 {
				Println("stor tiss")
				orders.ORDER_MATRIX[floor][i] = 0
				return 1
			}
		}
	}
	return 0
}

func Check_order(order_channel chan int) {

	for i := 0; i < N_FLOOR; i++ {
		for j := 0; j < N_BUTTON_TYPES; j++ {
			if orders.ORDER_MATRIX[i][j] > 0 {
				//Set_light()

				Elev_set_motor_direction(UP)
				Elev.Dir = UP
			}

		}
	}
}

func Event_manager(order_channel chan int, position_channel chan int) {
	for {

		select {
		//case order_received := <-order_channel:

		case current_floor := <-position_channel:
			Println("ballefrans")
			if Should_stop(order_channel, current_floor, Elev.Dir) != 0 {
				Stop_at_floor()
			}
			Check_order(order_channel)

		}
	}
}

func Stop_at_floor() {
	Elev_set_motor_direction(0)
	Elev_set_door_open_lamp(1)
	Sleep(2 * Second)
	Elev_set_door_open_lamp(0)
}
