package elevmanager

import (
	. "../driver"
	. "../message"
	. "../orders"
	. "fmt"
	"net"
)

type Elevator struct {
	Floor int
	Dir   Elev_dir
	//ORDER_INSIDE [N_FLOOR]int
}

type Elev_manager struct {
	Self_id         int
	All_elevators   map[int]*Elevator
	External_orders [2][N_FLOOR]int
	Master          int
}

func Make_elev_manager() Elev_manager {
	var elev Elev_manager
	elev.All_elevators = make(map[int]*Elevator)
	addr, _ := net.InterfaceAddrs()
	elev.Self_id = int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')
	elev.All_elevators[elev.Self_id] = new(Elevator)

	elev.All_elevators[elev.Self_id].Floor = Elev_get_floor_sensor_signal()
	elev.choose_master()
	return elev
}

func (elev *Elev_manager) choose_master() {
	current_min := 255
	for id := range elev.All_elevators {
		if id < current_min {
			current_min = id
		}
	}
	elev.Master = current_min
	Println("Master is ", elev.Master)
}

func (elev *Elev_manager) Set_elev_floor_and_direction(message UDPMessage) {
	_, ok := elev.All_elevators[message.Source]
	if ok {
		elev.All_elevators[message.Source].Floor = message.Floor
		elev.All_elevators[message.Source].Dir = message.Dir
	}
}

func (elev *Elev_manager) Add_elevator(message UDPMessage) { //might need to_network channel
	elev.All_elevators[message.Source] = new(Elevator)
	Println("Elevator ", message.Source, " is added")
	elev.choose_master()
}

func (elev *Elev_manager) Delete_elevator(id int) {
	delete(elev.All_elevators, id)
	elev.choose_master()
	Println("Elevator ", id, " is removed")
}

func (elev *Elev_manager) Assign_external_order(order External_order) int {
	elev_cost := make(map[int]int)
	for elevator, _ := range elev.All_elevators {
		elev_cost[elevator] = Calculate_cost(elev.All_elevators[elevator].Floor, elev.All_elevators[elevator].Dir, order)
		Println("Cost for ", elevator, ": ", elev_cost[elevator])
		Println("Elevator ", elevator, "is on floor", elev.All_elevators[elevator].Floor)
		Println("Elevator ", elevator, "has direction", elev.All_elevators[elevator].Dir)
		Println("")
	}
	best_elevator := -1
	min_cost := 1000
	for elevator, cost := range elev_cost {
		if cost < min_cost {
			min_cost = cost
			best_elevator = elevator
		}
	}
	Println("Ordertype ", order.Button_type, " in floor ", order.Floor, "assigned to elev", best_elevator)
	return best_elevator
}

func (elev *Elev_manager) Check_if_order_in_floor(message UDPMessage) bool {
	if elev.External_orders[message.Order.Button_type][message.Order.Floor] == 1 { /////kanskje denne skal kunne være 1/2 ut ifra status:knapp trykt, heis på vei
		return true
	}
	return false
}

func (elev *Elev_manager) Set_external_order(message UDPMessage) {
	elev.External_orders[message.Order.Button_type][message.Order.Floor] = 1
}
func (elev *Elev_manager) Clear_external_order(message UDPMessage) {
	if message.Order.Button_type == BUTTON_INSIDE {
		elev.External_orders[BUTTON_UP][message.Floor] = 0
		elev.External_orders[BUTTON_DOWN][message.Floor] = 0
	} else {
		elev.External_orders[message.Order.Button_type][message.Order.Floor] = 0
	}
}
