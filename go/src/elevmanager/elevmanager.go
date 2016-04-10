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

type elev_manager struct {
	self_id         int
	all_elevators   map[int]*Elevator
	external_orders [2][N_FLOOR]int
	master          int
}

func Make_elev_manager() elev_manager {
	var elev elev_manager
	elev.all_elevators = make(map[int]*Elevator)
	addr, _ := net.InterfaceAddrs()
	elev.self_id = int(addr[1].String()[12]-'0')*100 + int(addr[1].String()[13]-'0')*10 + int(addr[1].String()[14]-'0')
	elev.all_elevators[elev.self_id] = new(Elevator)

	elev.all_elevators[elev.self_id].Floor = Elev_get_floor_sensor_signal()
	elev.choose_master()
	return elev
}

func (elev *elev_manager) choose_master() {
	current_min := 255
	for id := range elev.all_elevators {
		if id < current_min {
			current_min = id
		}
	}
	elev.master = current_min
}

func (elev *elev_manager) Add_elevator(message UDPMessage) { //might need to_network channel
	elev.all_elevators[message.Source] = new(Elevator)
	Println("Elevator ", message.Source, " is added")
}

func (elev *elev_manager) Delete_elevator(id int, to_SM chan UDPMessage) {
	delete(elev.all_elevators, id)
	Println("Elevator ", id, " is removed")
}

func (elev *elev_manager) Assign_external_order(order External_order) int {
	elev_cost := make(map[int]int)
	for elevator, _ := range elev.all_elevators {
		elev_cost[elevator] = Calculate_cost(elev.all_elevators[elevator].Floor, elev.all_elevators[elevator].Dir, order)
	}
	best_elevator := -1
	min_cost := 1000
	for elevator, cost := range elev_cost {
		if cost < min_cost {
			min_cost = cost
			best_elevator = elevator
		}
	}
	return best_elevator
}
