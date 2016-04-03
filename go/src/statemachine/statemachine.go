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

var state ElevState = IDLE

type Elevator struct {
	Floor int
	Dir   Elev_dir
	Order int
	Order_type int
}

var Elev = Elevator{Elev_get_floor_sensor_signal(),STOPMOTOR,-1,-1}

func SM() {

	order_up_channel := make(chan int)
	order_down_channel := make(chan int)
	order_inside_channel := make(chan int)
	position_channel := make(chan int)
	order_channel:=make(chan int)
	timeout:=make(chan bool,1)

		
	go Check_current_order(order_channel)
	go Elevator_position(position_channel)
	go Check_order(order_up_channel, order_down_channel, order_inside_channel)
	go orders.Register_order()
	go orders.Print_orders()
	go Print_status()
	go Floor_timeout(timeout)
	
	Event_manager(/*order_up_channel, order_down_channel, order_inside_channel,*/order_channel,position_channel)
	

}



func Print_status(){
	for{
		Println("Current state: ",state)
		Println("Current direction: ", Elev.Dir)
		Println("Current floor: ", Elev.Floor)

		Sleep(1*Second)
	}
}

func Elevator_position(position_channel chan int) {
	for{
		floor := Elev_get_floor_sensor_signal()
		if floor != -1 {
			position_channel <- floor
			Elev_set_floor_indicator(floor)
		}
		Sleep(100 * Millisecond)
	}

}

func Check_current_order(order_channel chan int){
	for{
		Println("Current order: ",Elev.Order)
		if Elev.Order > -1{
			order_channel <- Elev.Order
		}
		Sleep(100*Millisecond)
	}
}

func Should_stop(floor int, dir Elev_dir) int {
	//Println("floor= ", floor)
	//Println("jeg har liten tiss")
	if orders.ORDER_INSIDE[floor] == 1 && Elev_get_floor_sensor_signal() >-1 {
				Println("flat tiss")
				//orders.ORDER_INSIDE[floor] = 0
				//orders.ORDER_UP[floor]=0
				//orders.ORDER_DOWN[floor] =0
				return 1
			}
	if floor == Elev.Order{
		Println("Hei")
		return 1
	}
	if dir == UP{
		Println("sannnn")
		if Elev.Order_type == BUTTON_UP || Elev.Order_type == BUTTON_INSIDE{
			if orders.ORDER_UP[floor] == 1{
				return 1
			}
		}
	}
	if dir == DOWN{
		Println("din mig")
		if Elev.Order_type == BUTTON_DOWN || Elev.Order_type == BUTTON_INSIDE{
			if orders.ORDER_DOWN[floor] == 1{
				return 1
			}
		}
	}
	
	return 0
}

func Check_order(order_up_channel, order_down_channel, order_inside_channel chan int) {
	for{
		//Println("Checking orders")
		floor:=Elev.Floor
		dir:=Elev.Dir
		if orders.No_orders() != 0{
			state=IDLE
			Elev.Dir=STOPMOTOR
		}
		switch state{
		case IDLE:
			Println("FUCK MY ASSSSSS")

			for i := 0; i < N_FLOOR; i++{

				if orders.ORDER_INSIDE[i] == 1{
					state = MOVING
					Elev.Order = i;
					Elev.Order_type = BUTTON_INSIDE
					//order_inside_channel <- i
				}
				if orders.ORDER_UP[i] == 1{
					state = MOVING
					Elev.Order = i;
					Elev.Order_type = BUTTON_UP
					//order_up_channel <- i
				}
				if orders.ORDER_DOWN[i] == 1{
					state = MOVING
					Elev.Order = i;
					Elev.Order_type = BUTTON_DOWN
					//order_down_channel <- i
				}
				
			}
		
		case MOVING:
			if dir == UP{
				Println("heidinkukksueger")

				if Elev.Order_type == BUTTON_UP || Elev.Order_type == BUTTON_INSIDE{
					if orders.ORDER_UP[floor] == 1 || orders.ORDER_INSIDE[floor] == 1{
						Elev.Order = floor
					}
				}
			}
				
			if dir == DOWN{
				if Elev.Order_type == BUTTON_DOWN || Elev.Order_type == BUTTON_INSIDE{
					if orders.ORDER_DOWN[floor] == 1 || orders.ORDER_INSIDE[floor] == 1{
						Elev.Order = floor
					}
				}
			}
			
		}


		
		Sleep(100 * Millisecond)
	}
}

func Event_manager(/*order_up_channel chan int,order_down_channel chan int,
	order_inside_channel chan int,*/order_channel chan int, position_channel chan int) {
	for {


		select {
		

		case current_order := <- order_channel:
			Println("Received Order")
			Move_to_floor(current_order)
		case current_floor := <-position_channel:
			Elev.Floor = current_floor
			if Should_stop(current_floor, Elev.Dir) != 0 {
				Stop_at_floor()
			}
			

		}
	}
}

func Stop_at_floor() {
	Println("Stopping at floor", Elev.Floor)
	orders.Clear_orders_at_floor(Elev.Floor)
	if Elev.Floor == Elev.Order{
		Elev.Order = -1
		state = IDLE
	} //else if Elev.Floor == Elev.Order{
	//	Elev.Order = -1
	//}
	orders.Clear_lights_at_floor(Elev.Floor)
	Elev_set_motor_direction(0)
	if Elev.Floor == 0{
		Elev.Dir = STOPMOTOR
	}
	if Elev.Floor == N_FLOOR -1{
		Elev.Dir = STOPMOTOR
	}
	Elev_set_door_open_lamp(1)
	Sleep(2 * Second)
	Elev_set_door_open_lamp(0)
}



func Move_to_floor(floor int){

	if floor > -1{
		Println("Moving to floor: ",floor)
		if floor < Elev.Floor{
			Elev_set_motor_direction(DOWN)
			Elev.Dir=DOWN
			state=MOVING
		}
		if floor > Elev.Floor{
			Elev_set_motor_direction(UP)
			Elev.Dir=UP
			state=MOVING
		}
		if floor == Elev.Floor && Elev_get_floor_sensor_signal() >-1{
			
			Stop_at_floor()
		}
	}
}

func Floor_timeout(timeout chan bool){
	for{
		Sleep(10 * Second)
		timeout <- true
	}

	Sleep(100 * Millisecond)
	
}
		