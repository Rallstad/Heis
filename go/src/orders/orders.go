package orders

import (
	. "../driver"
	. "fmt"
	. "time"
)

const N_ELEV = 3

type Queue struct {
	ORDER_UP     [N_FLOOR]int
	ORDER_DOWN   [N_FLOOR]int
	ORDER_INSIDE [N_FLOOR]int
}

type External_order struct {
	Floor       int
	Button_type int
}

var Elev_queue Queue

func Init_orders() {
	for i := 0; i < N_FLOOR; i++ {
		Elev_queue.ORDER_UP[i] = 0
		Elev_queue.ORDER_DOWN[i] = 0
		Elev_queue.ORDER_INSIDE[i] = 0
	}
	Elev_queue.ORDER_UP[N_FLOOR-1] = -1
	Elev_queue.ORDER_DOWN[0] = -1
}

func Register_order(order chan External_order) {
	for {
		Register_order_up(order)
		Register_order_down(order)
		Register_order_inside()

		Set_light()
		Sleep(100 * Millisecond)
	}

}

func Register_order_up(order chan External_order) {
	for i := 0; i < N_FLOOR-1; i++ {
		if Elev_get_button_signal(BUTTON_UP, i) > 0 {
			order <- External_order{Floor: i, Button_type: BUTTON_UP}
		}
	}

}

func Register_order_down(order chan External_order) {
	for i := 1; i < N_FLOOR; i++ {
		if Elev_get_button_signal(BUTTON_DOWN, i) > 0 {
			order <- External_order{Floor: i, Button_type: BUTTON_DOWN}
		}
	}

}

func Register_order_inside() {
	for i := 0; i < N_FLOOR; i++ {
		if Elev_get_button_signal(BUTTON_INSIDE, i) > 0 {
			Elev_queue.ORDER_INSIDE[i] = 1
		}
	}

}

func Place_order(order External_order) {
	if order.Button_type == BUTTON_UP {
		Elev_queue.ORDER_UP[order.Floor] = 1
	}
	if order.Button_type == BUTTON_DOWN {
		Elev_queue.ORDER_DOWN[order.Floor] = 1
	}
}

func Calculate_cost(elev_pos int, elev_dir Elev_dir, order External_order) int {
	cost := 0
	order_dir := elev_pos - order.Floor
	if order_dir*int(elev_dir) > 0 {
		cost += 10
	} else if order_dir*int(elev_dir) < 0 {
		if elev_dir == UP && order.Button_type == BUTTON_DOWN {
			cost += 3
		} else if elev_dir == DOWN && order.Button_type == BUTTON_UP {
			cost += 3
		}
	}
	return cost
}

func Clear_orders_at_floor(floor int) {
	if floor == 0 {
		Elev_queue.ORDER_UP[floor] = 0
		Elev_queue.ORDER_INSIDE[floor] = 0
	} else if floor == N_FLOOR-1 {
		Elev_queue.ORDER_DOWN[floor] = 0
		Elev_queue.ORDER_INSIDE[floor] = 0
	} else {
		Elev_queue.ORDER_UP[floor] = 0
		Elev_queue.ORDER_DOWN[floor] = 0
		Elev_queue.ORDER_INSIDE[floor] = 0
	}
}

func Clear_lights_at_floor(floor int) {
	if floor == 0 {
		Elev_set_button_lamp(BUTTON_UP, floor, 0)
		Elev_set_button_lamp(BUTTON_INSIDE, floor, 0)
	} else if floor == N_FLOOR-1 {
		Elev_set_button_lamp(BUTTON_DOWN, floor, 0)
		Elev_set_button_lamp(BUTTON_INSIDE, floor, 0)
	} else {
		Elev_set_button_lamp(BUTTON_UP, floor, 0)
		Elev_set_button_lamp(BUTTON_DOWN, floor, 0)
		Elev_set_button_lamp(BUTTON_INSIDE, floor, 0)
	}
}

func No_orders() int {
	for i := 0; i < N_FLOOR; i++ {
		if Elev_queue.ORDER_UP[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_DOWN[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_INSIDE[i] == 1 {
			return 0
		}
	}
	//	Println("NO MORE ORDERS")
	return 1
}

func No_orders_above(floor int) int {
	for i := floor; i < N_FLOOR; i++ {
		if Elev_queue.ORDER_UP[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_DOWN[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_INSIDE[i] == 1 {
			return 0
		}
	}
	return 1
}

func No_orders_below(floor int) int {
	for i := floor; i > -1; i-- {
		if Elev_queue.ORDER_UP[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_DOWN[i] == 1 {
			return 0
		}
		if Elev_queue.ORDER_INSIDE[i] == 1 {
			return 0
		}
	}
	return 1
}

func No_orders_inside() int {
	for i := 0; i < N_FLOOR; i++ {
		if Elev_queue.ORDER_INSIDE[i] > 0 {
			return 0
		}
	}
	return 1
}

func Set_light() {
	for i := 0; i < N_FLOOR; i++ {
		if Elev_queue.ORDER_UP[i] > 0 {
			Elev_set_button_lamp(BUTTON_UP, i, 1)
		}
		if Elev_queue.ORDER_DOWN[i] > 0 {
			Elev_set_button_lamp(BUTTON_DOWN, i, 1)
		}
		if Elev_queue.ORDER_INSIDE[i] > 0 {
			Elev_set_button_lamp(BUTTON_INSIDE, i, 1)
		}
	}
}

func Print_orders() {
	for {
		Println("Current orders")
		for i := 0; i < N_FLOOR; i++ {
			Println(Elev_queue.ORDER_UP[i])
			Println(Elev_queue.ORDER_DOWN[i])
			Println(Elev_queue.ORDER_INSIDE[i])
		}
		Sleep(3 * Second)
	}
}
