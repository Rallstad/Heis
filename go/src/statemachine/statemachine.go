package statemachine

import (
	. "../driver"
	"../orders"
	. "fmt"
	. "time"
)

type ElevState int

const (
	//INIT ElevState = iota
	IDLE ElevState = iota
	MOVING
	STOPPED
	DOOR_OPEN
)

const (
	stop      = 0
	move_up   = 1
	move_down = -1
)

var state ElevState = IDLE

type Elevator struct {
	Floor int
	Dir   Elev_dir
}

var Elev = Elevator{Elev_get_floor_sensor_signal(), STOPMOTOR}

func SM() {

	position_channel := make(chan int)
	order_channel := make(chan int)
	command_channel := make(chan int)

	go Elevator_position(position_channel)
	go Check_order(order_channel)
	go orders.Register_order()
	go orders.Print_orders()
	go Print_status()

	go Command_manager(order_channel, position_channel, command_channel)
	Event_manager(order_channel, position_channel, command_channel)

}

func Print_status() {
	for {
		Println("Current state: ", state)
		Println("Current direction: ", Elev.Dir)
		Println("Current floor: ", Elev.Floor)

		Sleep(1 * Second)
	}
}

func Elevator_position(position_channel chan int) {
	for {
		floor := Elev_get_floor_sensor_signal()
		if floor != -1 {
			position_channel <- floor
			Elev_set_floor_indicator(floor)
		}
		Sleep(100 * Millisecond)
	}

}

func Should_stop(floor int, dir Elev_dir, command_channel chan int) {
	Println("Checking if stop")
	if orders.Elev_queue.ORDER_INSIDE[floor] == 1 && Elev_get_floor_sensor_signal() > -1 {
		if floor == 0 || floor == N_FLOOR-1 {
			state = IDLE
		}
		Println("Should stop order inside")
		command_channel <- stop
		Stop_at_floor()
	} else if dir == UP {
		Println("Saggy tits")
		if orders.Elev_queue.ORDER_UP[floor] == 1 {
			Println("Stopping for order UP")
			command_channel <- stop
			Stop_at_floor()
		} else if floor == N_FLOOR-1 { ///KANSKJE FJERNES
			Println("Stopping for top floor")
			//state = IDLE
			command_channel <- stop
			Stop_at_floor()
		} else if orders.Elev_queue.ORDER_DOWN[floor] == 1 && orders.No_orders_above(floor+1) != 0 {
			command_channel <- stop
			Stop_at_floor()
		}
	} else if dir == DOWN {
		Println("stiff niples")
		if orders.Elev_queue.ORDER_DOWN[floor] == 1 {
			Println("Stopping for order DOWN")
			command_channel <- stop
			Stop_at_floor()
		} else if floor == 0 { ///KANSKJE FJERNES
			Println("Stopping for bottom floor")
			//state = IDLE
			command_channel <- stop
			Stop_at_floor()
		} else if orders.Elev_queue.ORDER_UP[floor] == 1 && orders.No_orders_below(floor-1) != 0 {
			command_channel <- stop
			Stop_at_floor()
		}
	}
}

func Get_next_direction(command_channel chan int) {
	if Elev.Dir == UP {
		if orders.No_orders_above(Elev.Floor) == 0 {
			command_channel <- move_up
		} else {
			Elev.Dir = STOPMOTOR
			command_channel <- stop
		}

	} else if Elev.Dir == DOWN {
		if orders.No_orders_below(Elev.Floor) == 0 {
			command_channel <- move_down
		} else {
			Elev.Dir = STOPMOTOR
			command_channel <- stop
		}
	} else if Elev.Dir == STOPMOTOR {
		for i := Elev.Floor + 1; i < N_FLOOR; i++ {
			if orders.Elev_queue.ORDER_UP[i] == 1 || orders.Elev_queue.ORDER_DOWN[i] == 1 || orders.Elev_queue.ORDER_INSIDE[i] == 1 {
				command_channel <- move_up
			}
		}
		for i := 0; i < Elev.Floor; i++ {
			if orders.Elev_queue.ORDER_UP[i] == 1 || orders.Elev_queue.ORDER_DOWN[i] == 1 || orders.Elev_queue.ORDER_INSIDE[i] == 1 {
				command_channel <- move_down
			}
		}
		if orders.Elev_queue.ORDER_UP[Elev.Floor] == 1 || orders.Elev_queue.ORDER_DOWN[Elev.Floor] == 1 || orders.Elev_queue.ORDER_INSIDE[Elev.Floor] == 1 {
			command_channel <- stop
		}
	}
}

func Check_order(order_channel chan int) {
	for {

		//dir := Elev.Dir
		if orders.No_orders() != 0 {
			state = IDLE
			Elev.Dir = STOPMOTOR
		}
		switch state {
		case IDLE:
			Println("case idle i check order")

			for i := 0; i < N_FLOOR; i++ {

				if orders.Elev_queue.ORDER_INSIDE[i] == 1 {
					order_channel <- i

				}
				if orders.Elev_queue.ORDER_UP[i] == 1 {
					order_channel <- i

				}
				if orders.Elev_queue.ORDER_DOWN[i] == 1 {
					order_channel <- i
				}

			}

			/*case MOVING:
			if dir == UP {
				Println("moving dir up")

				for i := Elev.Floor; i < N_FLOOR; i++ {
					if orders.Elev_queue.ORDER_UP[i] == 1 {
						order_channel <- i
					} else if orders.Elev_queue.ORDER_INSIDE[i] == 1 {
						order_channel <- i
					} else if orders.Elev_queue.ORDER_DOWN[i] == 1 {
						order_channel <- i
					}
				}
			}
			if dir == DOWN {
				Println("moving dir down")
				for i := Elev.Floor; i <= 0; i-- {
					if orders.Elev_queue.ORDER_UP[i] == 1 {
						order_channel <- i
					} else if orders.Elev_queue.ORDER_INSIDE[i] == 1 {
						order_channel <- i
					} else if orders.Elev_queue.ORDER_DOWN[i] == 1 {
						order_channel <- i
					}
				}
			}*/
		}
		Sleep(100 * Millisecond)
	}
}

func Command_manager(order_channel chan int, position_channel chan int, command_channel chan int) {
	for {

		select {
		case command := <-command_channel:
			switch command {
			case stop:
				Println("Received STOP command")
				Elev_set_motor_direction(STOPMOTOR)
				//Stop_at_floor()
				break
			case move_up:
				Println("Received MOVEUP command")
				Elev_set_motor_direction(UP)
				Elev.Dir = UP
				break
			case move_down:
				Println("Received MOVEDOWN command")
				Elev_set_motor_direction(DOWN)
				Elev.Dir = DOWN
				break
			}
		}
	}
}

func Event_manager(order_channel chan int, position_channel chan int, command_channel chan int) {
	for {
		select {
		case current_floor := <-position_channel:
			Elev.Floor = current_floor
			Should_stop(current_floor, Elev.Dir, command_channel)
			Get_next_direction(command_channel)

		}
	}
}

func Stop_at_floor() {
	Println("Stopping at floor", Elev.Floor)
	orders.Clear_orders_at_floor(Elev.Floor)
	orders.Clear_lights_at_floor(Elev.Floor)
	Elev_set_motor_direction(STOPMOTOR)

	if Elev.Floor == 0 {
		Elev.Dir = UP
	}
	if Elev.Floor == N_FLOOR-1 {
		Elev.Dir = DOWN
	}
	Elev_set_door_open_lamp(1)
	Sleep(2 * Second)
	Elev_set_door_open_lamp(0)
}

func Move_to_floor(floor int) {

	if floor > -1 {
		Println("Moving to floor: ", floor)
		if floor < Elev.Floor {
			Elev_set_motor_direction(DOWN)
			Elev.Dir = DOWN
			state = MOVING
		}
		if floor > Elev.Floor {
			Elev_set_motor_direction(UP)
			Elev.Dir = UP
			state = MOVING
		}
		if floor == Elev.Floor && Elev_get_floor_sensor_signal() > -1 {

			Stop_at_floor()
		}
	}
}
