package statemachine

import (
	. "../driver"
	. "../elevmanager"
	. "../message"
	. "../network"
	"../orders"
	. "fmt"
	"os"
	"os/signal"
	"syscall"
	. "time"
)

type ElevState int

const (
	//INIT ElevState = iota
	IDLE ElevState = iota
	MOVING
	DOOR_OPEN
)

const (
	stop      = 0
	move_up   = 1
	move_down = -1
)

var state ElevState = IDLE

var Elev = Elevator{Elev_get_floor_sensor_signal(), STOPMOTOR}

func SM() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	signal.Notify(exit, syscall.SIGTERM)
	from_SM := make(chan UDPMessage, 1000)
	to_SM := make(chan UDPMessage, 1000)

	position_channel := make(chan int, 100)
	order_channel := make(chan int, 100)
	ext_order_channel := make(chan orders.External_order, 100)
	command_channel := make(chan int, 100)

	go Elevator_position(position_channel)
	//go Check_order(order_channel)
	go orders.Register_order(ext_order_channel)
	//go orders.Print_orders()
	//go Print_status()
	go Network_manager(from_SM, to_SM)
	go Command_manager(command_channel, from_SM)

	go func() {
		<-exit
		handle_program_exit()
		os.Exit(1)
	}()
	Event_manager(ext_order_channel, order_channel, position_channel, command_channel, from_SM, to_SM)

}

func handle_program_exit() {
	Println(" Program stopped by human")
	Elev_set_motor_direction(STOPMOTOR)
}

func Print_status() {
	for {
		Println("Current state: ", state)
		Println("Current direction: ", Elev.Dir)
		Println("Current floor: ", Elev.Floor)
		Sleep(2 * Second)
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

func Should_stop(floor int, dir Elev_dir, command_channel chan int, from_SM chan UDPMessage) {
	from_SM <- UDPMessage{Floor: floor}
	//Println("Checking if stop")
	if orders.Elev_queue.ORDER_INSIDE[floor] == 1 && Elev_get_floor_sensor_signal() > -1 {
		if floor == 0 || floor == N_FLOOR-1 {
			state = IDLE
			//Elev.Dir = STOPMOTOR
		}
		//Println("Should stop order inside")
		command_channel <- stop
		from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
		Stop_at_floor()
	} else if dir == UP {
		//Println("Saggy tits")
		if orders.Elev_queue.ORDER_UP[floor] == 1 {
			Println("Stopping in floor ", floor)
			//Println("Stopping for order UP")
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		} else if floor == N_FLOOR-1 { ///KANSKJE FJERNES
			//Println("Stopping for top floor")
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		} else if orders.Elev_queue.ORDER_DOWN[floor] == 1 && orders.No_orders_above(floor+1) != 0 {
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		}
	} else if dir == DOWN {
		//Println("stiff tits")
		if orders.Elev_queue.ORDER_DOWN[floor] == 1 {
			//Println("Stopping for order DOWN")
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		} else if floor == 0 { ///KANSKJE FJERNES
			//Println("Stopping for bottom floor")
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		} else if orders.Elev_queue.ORDER_UP[floor] == 1 && orders.No_orders_below(floor-1) != 0 {
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		}
	} else if dir == STOPMOTOR {
		//Println("Fat ass")
		if orders.Elev_queue.ORDER_UP[floor] == 1 || orders.Elev_queue.ORDER_DOWN[floor] == 1 {
			command_channel <- stop
			from_SM <- UDPMessage{MessageId: Order_completed, Floor: floor, Dir: dir}
			Stop_at_floor()
		}
	}

}

func Get_next_direction(command_channel chan int) {
	//Println("Getting next direction")
	if Elev.Dir == UP {
		if orders.No_orders_above(Elev.Floor) == 0 {
			command_channel <- move_up
		} else {
			//Elev.Dir = STOPMOTOR
			command_channel <- stop
		}

	} else if Elev.Dir == DOWN {
		if orders.No_orders_below(Elev.Floor) == 0 {
			command_channel <- move_down
		} else {
			//Elev.Dir = STOPMOTOR
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

func Command_manager(command_channel chan int, from_SM chan UDPMessage) {
	for {

		select {
		case command := <-command_channel:
			switch command {
			case stop:
				//Println("Received STOP command")
				Elev_set_motor_direction(STOPMOTOR)
				Elev.Dir = STOPMOTOR
				from_SM <- UDPMessage{MessageId: Elev_state_update, Floor: Elev.Floor, Dir: STOPMOTOR}
				break
			case move_up:
				//Println("Received MOVEUP command")
				Elev_set_motor_direction(UP)
				Elev.Dir = UP
				from_SM <- UDPMessage{MessageId: Elev_state_update, Floor: Elev.Floor, Dir: UP}
				break
			case move_down:
				//Println("Received MOVEDOWN command")
				Elev_set_motor_direction(DOWN)
				Elev.Dir = DOWN
				from_SM <- UDPMessage{MessageId: Elev_state_update, Floor: Elev.Floor, Dir: DOWN}
				break
			}
		}
	}
}

func Event_manager(ext_order_channel chan orders.External_order, order_channel chan int, position_channel chan int, command_channel chan int, from_SM chan UDPMessage, to_SM chan UDPMessage) {
	elev := Make_elev_manager()
	for {
		select {
		case message := <-to_SM:
			switch message.MessageId {
			case Elev_state_update:
				elev.Set_elev_floor_and_direction(message)
			case Elev_add:
				elev.Add_elevator(message)
				from_SM <- UDPMessage{MessageId: Elev_state_update, Floor: Elev.Floor}
				break
			case Elev_delete:
				elev.Delete_elevator(message.Source)
			case New_order:
				//Println("ext_butt_pushed_New_order")
				if !elev.Check_if_order_in_floor(message) {
					elev.Set_external_order(message)
					orders.Set_ext_light(message.Order)
					if elev.Self_id == elev.Master {
						from_SM <- UDPMessage{MessageId: Order_assigned, Target: elev.Assign_external_order(message.Order), Order: message.Order, Floor: Elev.Floor}
					}
				}
				break
			case Order_assigned:
				//Println("ord ass")
				if message.Target == elev.Self_id {
					orders.Place_order(message.Order)
				}
				break
			case Order_completed:
				elev.Clear_external_order(message)
				orders.Clear_ext_light(message.Floor)
				break
			}
		case current_floor := <-position_channel:
			//Println("curr fl")
			Elev.Floor = current_floor
			Should_stop(current_floor, Elev.Dir, command_channel, from_SM)
			Get_next_direction(command_channel)

		case ext_order := <-ext_order_channel:
			Println("ext_order")
			from_SM <- UDPMessage{MessageId: New_order, Order: ext_order}
		}
	}
}

func Stop_at_floor() {
	//Println("Stopping at floor", Elev.Floor)
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
