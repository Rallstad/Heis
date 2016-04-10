package message

import "../orders"

//Message ID
const (
	Ping            = 1
	Elev_move       = 2
	New_order       = 3
	Order_assigned  = 4
	Order_completed = 5
	Elev_delete     = 6
	Elev_add        = 7
)

type UDPMessage struct {
	MessageId int
	Source    int
	Target    int
	Floor     int
	State     [2]int
	Order     orders.External_order
}

/*func CalculateChecksum(Msg *UDPMessage) int { // not a very good crc, just for testing
	c := Msg.MessageId%7 + Msg.OrderQueue[0]%7 + Msg.ElevatorStateUpdate[0]%7
	return c
}*/
