package orders

import (
	. "../driver"
	. "fmt"
	. "time"
)

const N_ELEV = 3

var ORDER_UP =make([]int,N_FLOOR)
var ORDER_DOWN =make([]int,N_FLOOR)
var ORDER_INSIDE =make([]int,N_FLOOR)

func Init_orders(){
	for i:=0;i<N_FLOOR;i++{
		ORDER_UP[i]=0
		ORDER_DOWN[i]=0
		ORDER_INSIDE[i]=0
	}
	ORDER_UP[N_FLOOR-1]=-1
	ORDER_DOWN[0]=-1
}

func Register_order() {
	for{
		Register_order_up()
		Register_order_down()
		Register_order_inside()
		
		Set_light()
		Sleep(100*Millisecond)
	}

}

func Register_order_up(){
	for i := 0; i < N_FLOOR-1; i++{
		if Elev_get_button_signal(BUTTON_UP,i) >0{
			ORDER_UP[i]=1
		}
	}

}

func Register_order_down(){
	for i := 1; i < N_FLOOR; i++{
		if Elev_get_button_signal(BUTTON_DOWN,i) >0{
			ORDER_DOWN[i]=1
		}
	}

}

func Register_order_inside(){
	for i := 0; i < N_FLOOR; i++{
		if Elev_get_button_signal(BUTTON_INSIDE,i) >0{
			ORDER_INSIDE[i]=1
		}
	}

}

func Clear_orders_at_floor(floor int){
	if floor == 0{
		ORDER_UP[floor] = 0
		ORDER_INSIDE[floor] = 0
	} else if floor == N_FLOOR -1{
		ORDER_DOWN[floor] = 0
		ORDER_INSIDE[floor] = 0
	} else{
		ORDER_UP[floor] = 0
		ORDER_DOWN[floor] = 0
		ORDER_INSIDE[floor] = 0
	}
}

func Clear_lights_at_floor(floor int){
	if floor == 0{
		Elev_set_button_lamp(BUTTON_UP,floor,0)
		Elev_set_button_lamp(BUTTON_INSIDE,floor,0)
	} else if floor == N_FLOOR -1{
		Elev_set_button_lamp(BUTTON_DOWN,floor,0)
		Elev_set_button_lamp(BUTTON_INSIDE,floor,0)
	} else{
		Elev_set_button_lamp(BUTTON_UP,floor,0)
		Elev_set_button_lamp(BUTTON_DOWN,floor,0)
		Elev_set_button_lamp(BUTTON_INSIDE,floor,0)
	}
}

func No_orders() int {
	for i := 0; i < N_FLOOR; i++{
		if ORDER_UP[i]==1{
			return 0
		}
		if ORDER_DOWN[i]==1{
			return 0
		}
		if ORDER_INSIDE[i]==1{
			return 0
		}
	}
	Println("NO MORE ORDERS")
	return 1
}

func No_orders_above(floor int) int{
	for i := floor; i < N_FLOOR; i++{
		if ORDER_UP[i]==1{
			return 0
		}
		if ORDER_DOWN[i]==1{
			return 0
		}
		if ORDER_INSIDE[i]==1{
			return 0
		}
	}
	return 1
}


func No_orders_below(floor int) int{
	for i := floor; i > -1; i-- {
		if ORDER_UP[i]==1{
			return 0
		}
		if ORDER_DOWN[i]==1{
			return 0
		}
		if ORDER_INSIDE[i]==1{
			return 0
		}
	}
	return 1
}

func No_orders_inside() int {
	for i := 0; i < N_FLOOR; i++{
		if ORDER_INSIDE[i] > 0{
			return 0
		}
	}
	return 1
}


func Set_light() {
	for i := 0; i < N_FLOOR; i++ {
		if ORDER_UP[i] > 0 {
			Elev_set_button_lamp(BUTTON_UP, i, 1)
		}
		if ORDER_DOWN[i] > 0 {
			Elev_set_button_lamp(BUTTON_DOWN, i, 1)
		}
		if ORDER_INSIDE[i] > 0 {
			Elev_set_button_lamp(BUTTON_INSIDE, i, 1)
		}
	}
}

func Print_orders() {
	for{
	Println("Current orders")
	for i := 0; i < N_FLOOR; i++ {
		Println(ORDER_UP[i])
		Println(ORDER_DOWN[i])
		Println(ORDER_INSIDE[i])
		}
		Sleep(3*Second)
	}
}
