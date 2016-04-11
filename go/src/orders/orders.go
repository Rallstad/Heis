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
			Println("SOMEONE PRESSED THE ORDER UP BUTTON IN FLOOR ", i, "!!!!!!!!!!!!!")
			Sleep(100 * Millisecond)
		}
	}

}

func Register_order_down(order chan External_order) {
	for i := 1; i < N_FLOOR; i++ {
		if Elev_get_button_signal(BUTTON_DOWN, i) > 0 {
			order <- External_order{Floor: i, Button_type: BUTTON_DOWN}
			Println("SOMEONE PRESSED THE ORDER DOWN BUTTON IN FLOOR ", i, "!!!!!!!!!!!!!")
			Sleep(100 * Millisecond)
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
	Println("Placing order ", order.Button_type, " in floor ", order.Floor)
	if order.Button_type == BUTTON_UP {
		Elev_queue.ORDER_UP[order.Floor] = 1
	}
	if order.Button_type == BUTTON_DOWN {
		Elev_queue.ORDER_DOWN[order.Floor] = 1
	}
}

func Calculate_cost(elev_pos int, elev_dir Elev_dir, order External_order) int {

	cost := 0
	if elev_pos == order.Floor && elev_dir == STOPMOTOR {
		return cost
	}
	if elev_dir != STOPMOTOR {
		cost += 4
	}
	order_dir := elev_pos - order.Floor
	if order_dir < 0 {
		cost += (-order_dir)
	} else if order_dir > 0 {
		cost += order_dir
	}
	if order_dir*int(elev_dir) > 0 {
		Println("Elev in floor ", elev_pos, " has to switch direction")
		cost += 10
	} else if order_dir*int(elev_dir) < 0 {
		if elev_dir == UP && order.Button_type == BUTTON_DOWN {
			cost += 5
		} else if elev_dir == DOWN && order.Button_type == BUTTON_UP {
			cost += 5
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
	Println("Clearing lights in floor ", floor)
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

func Set_ext_light(order External_order) {
	Elev_set_button_lamp(order.Button_type, order.Floor, 1)
}

func Clear_ext_light(floor int) {
	Println("Clearing external light in floor", floor) //, " ,button: ", order.Button_type)
	Elev_set_button_lamp(BUTTON_UP, floor, 0)
	Elev_set_button_lamp(BUTTON_DOWN, floor, 0)
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
